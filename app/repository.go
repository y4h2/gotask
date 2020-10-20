package app

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/y4h2/gotask/internal/dao"
)

type TaskRepository struct {
	db *sqlx.DB
}

func (*TaskRepository) TableName() string {
	return "gotask.task"
}

func (repo *TaskRepository) List() ([]dao.Task, error) {
	var tasks = []dao.Task{}

	err := repo.db.Select(&tasks, fmt.Sprintf(`
		SELECT id, type, host FROM %s`, repo.TableName()))
	return tasks, err
}

func (repo *TaskRepository) Read(id string) (dao.Task, error) {
	var task = dao.Task{}

	err := repo.db.Get(&task, fmt.Sprintf(
		"SELECT id, type, host FROM %s WHERE id=$1", repo.TableName()),
		id)
	return task, err
}

func (repo *TaskRepository) ReadOrCreate(id string, host string, typ string) (task dao.Task, exist bool, err error) {
	task = dao.Task{
		ID:   id,
		Type: typ,
	}
	err = repo.db.QueryRow(fmt.Sprintf(`
	WITH e AS(
			INSERT INTO %s ("id", "host", "type") 
						VALUES ($1, $2, $3)
			ON CONFLICT("id") DO NOTHING
			RETURNING host, false
	)
	SELECT * FROM e
	UNION
			SELECT host, true FROM %s WHERE id=$1`, repo.TableName(), repo.TableName()),
		id, host, typ).
		Scan(&task.Host, &exist)

	return
}

func (repo *TaskRepository) Delete(id string) (bool, error) {
	result, err := repo.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=$1", repo.TableName()), id)
	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affectedRows >= 1, nil
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}
