package app

import "github.com/y4h2/gotask/internal/dao"

type Task struct {
	ID   string
	Host string
	Type string
}

func DaoToEntityTask(t dao.Task) Task {
	return Task{
		ID:   t.ID,
		Host: t.Host,
		Type: t.Type,
	}
}
