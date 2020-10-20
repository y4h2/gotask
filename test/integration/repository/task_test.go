package task_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/y4h2/gotask/app"
	"github.com/y4h2/gotask/internal/dao"
	"github.com/y4h2/gotask/internal/enum"
	"github.com/y4h2/gotask/test/testutil"
)

type TaskRepositoryTestSuite struct {
	testutil.Suite

	db   *sqlx.DB
	repo *app.TaskRepository
}

func (suite *TaskRepositoryTestSuite) SetupSuite() {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("user=myuser port=%d password=mypassword sslmode=disable", 5432))
	if err != nil {
		suite.Require().NoError(err)
	}
	suite.db = db

	suite.repo = app.NewTaskRepository(suite.db)
}

func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *TaskRepositoryTestSuite) SetupTest() {

}

func (suite *TaskRepositoryTestSuite) TearDownTest() {
	suite.cleanup()
}

func (suite *TaskRepositoryTestSuite) newTask() dao.Task {
	return dao.Task{
		ID:   gofakeit.UUID(),
		Host: gofakeit.IPv4Address(),
		Type: gofakeit.RandomString([]string{enum.TaskType.PRINT, enum.TaskType.SLEEP}),
	}
}

func (suite *TaskRepositoryTestSuite) createTask() dao.Task {
	task := suite.newTask()
	_, err := suite.db.NamedExec(fmt.Sprintf(`
		INSERT INTO %s (id, host, type)
		VALUES (:id, :host, :type)`, suite.repo.TableName()),
		task)
	suite.Require().NoError(err)

	return task
}

func (suite *TaskRepositoryTestSuite) removeTask(id string) {
	suite.db.MustExec(
		fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, suite.repo.TableName()),
		id)
}

func (suite *TaskRepositoryTestSuite) listTasks() []dao.Task {
	var tasks = []dao.Task{}

	suite.Require().NoError(
		suite.db.Select(&tasks, fmt.Sprintf(`
		SELECT id, type, host FROM %s`, suite.repo.TableName())))

	return tasks
}

func (suite *TaskRepositoryTestSuite) cleanup() {
	suite.db.MustExec("DELETE FROM " + suite.repo.TableName())
}

func (suite *TaskRepositoryTestSuite) TestList() {
	t := suite.T()
	t.Log("Given a task")
	task := suite.createTask()

	t.Log("When I call the List function")
	tasks, err := suite.repo.List()
	suite.NoError(err)

	t.Log("Then I should get the given task")
	suite.Equal(1, len(tasks))
	suite.Equal(task, tasks[0])
}

func (suite *TaskRepositoryTestSuite) TestRead() {
	suite.Given("a task")
	task := suite.createTask()

	suite.When("I call the Read function")
	returnedTask, err := suite.repo.Read(task.ID)
	suite.NoError(err)

	suite.Then("I should get the given task")
	suite.Equal(returnedTask, task)
}

func (suite *TaskRepositoryTestSuite) TestReadOrCreate() {
	suite.When("I call the ReadOrCreate function directly")
	task := suite.newTask()
	returnedTask, exist, err := suite.repo.ReadOrCreate(task.ID, task.Host, task.Type)
	suite.NoError(err)

	suite.Then("it should show the data doesn't exist")
	suite.False(exist)

	suite.Then("the returned task should match with out input")
	suite.Equal(returnedTask, task)

	suite.When("I call the ReadOrCreate function with the same task again")
	returnedTask, exist, err = suite.repo.ReadOrCreate(task.ID, task.Host, task.Type)
	suite.NoError(err)

	suite.Then("it should show the data exist")
	suite.True(exist)

	suite.Then("the returned task should match with out input")
	suite.Equal(returnedTask, task)
}

func (suite *TaskRepositoryTestSuite) TestDelete() {
	suite.Given("a task")
	task := suite.createTask()

	suite.When("I call the Read function")
	deleted, err := suite.repo.Delete(task.ID)
	suite.NoError(err)
	suite.True(deleted)

	suite.Then("the db should have no data")
	tasks := suite.listTasks()
	suite.Equal(0, len(tasks))
}

func TestTaskRepository(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
