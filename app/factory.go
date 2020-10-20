package app

import (
	"errors"
	"fmt"
	"time"
)

type TaskRunner interface {
	Run(stopChan chan struct{}) error
}

func BuildTask(taskType string) (TaskRunner, error) {
	if taskType == "print" {
		return NewPrintTask(), nil
	}

	return nil, errors.New("unknow task type")
}

type PrintTask struct {
}

func (PrintTask) Run(stopChan chan struct{}) error {
	for cnt := 0; cnt < 1000; cnt++ {
		select {
		case <-stopChan:
			return nil
		default:
			fmt.Println("task")
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func NewPrintTask() *PrintTask {
	return &PrintTask{}
}
