package app

import (
	"golang.org/x/sync/errgroup"
)

type TaskFactory interface {
	GetTaskRunner(taskType string) (TaskRunner, error)
}

type Notifier interface {
	NotifySuccess() error
	NotifyFail() error
}

type TaskManager struct {
	controller   map[string]chan struct{}
	taskFactory  TaskFactory
	taskNotifier Notifier
}

func (mgr *TaskManager) IsTaskRunning(taskID string) bool {
	_, ok := mgr.controller[taskID]
	return ok
}

func (mgr *TaskManager) CreateTask(taskID string, taskType string) error {
	stopChan := make(chan struct{})
	mgr.controller[taskID] = stopChan

	task, err := mgr.taskFactory.GetTaskRunner(taskType)
	if err != nil {
		return err
	}

	var group errgroup.Group
	group.Go(func() error {
		return task.Run(stopChan)
	})

	if err := group.Wait(); err != nil {
		mgr.taskNotifier.NotifyFail()
	}
	mgr.taskNotifier.NotifySuccess()

	return nil
}

func (mgr *TaskManager) CancelTask(taskID string) error {
	stopChan, ok := mgr.controller[taskID]
	if !ok {
		return nil
	}
	close(stopChan)

	delete(mgr.controller, taskID)

	return nil
}

func NewTaskManager(taskFactory TaskFactory, taskNotifier Notifier) *TaskManager {
	return &TaskManager{
		controller:   make(map[string]chan struct{}),
		taskFactory:  taskFactory,
		taskNotifier: taskNotifier,
	}
}
