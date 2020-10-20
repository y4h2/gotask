package app

import "github.com/y4h2/gotask/internal/dao"

type taskRepository interface {
	List() ([]dao.Task, error)
	Read(name string) (dao.Task, error)
	ReadOrCreate(name string, host string) (task dao.Task, exist bool, err error)
	Delete(name string) (bool, error)
}

type taskManager interface {
	IsTaskRunning(taskID string) bool
	CreateTask(taskID string) error
	CancelTask(taskID string) error
}

type interCallAdapter interface {
	SendCancelRequest(host string, taskID string) error
}

type TaskService struct {
	repository       taskRepository
	taskManager      taskManager
	interCallAdapter interCallAdapter
	host             string
}

func (srv *TaskService) List() ([]Task, error) {
	daoTasks, err := srv.repository.List()
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, len(daoTasks))
	for i := range daoTasks {
		tasks[i] = DaoToEntityTask(daoTasks[i])
	}

	return tasks, nil
}

func (srv *TaskService) GetByName(taskID string) (Task, error) {
	task, err := srv.repository.Read(taskID)
	if err != nil {
		return Task{}, err
	}
	return DaoToEntityTask(task), nil
}

func (srv *TaskService) Create(taskID string) {
	srv.repository.ReadOrCreate(taskID, srv.host)
	srv.taskManager.CreateTask(taskID)

}

func (srv *TaskService) Cancel(taskID string) error {
	if srv.taskManager.IsTaskRunning(taskID) {
		err := srv.taskManager.CancelTask(taskID)
		if err != nil {
			return err
		}
		srv.repository.Delete(taskID)

		return nil
	}
	// redirect
	task, err := srv.repository.Read(taskID)
	if err != nil {
		return err
	}

	return srv.interCallAdapter.SendCancelRequest(task.Host, taskID)
}

func NewTaskService(repository taskRepository, host string) *TaskService {
	return &TaskService{
		repository: repository,
		host:       host,
	}
}
