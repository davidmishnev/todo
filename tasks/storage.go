package tasks

import (
	"math/rand"
	"sync"
)

type Storage struct {
	mutex sync.RWMutex
	tasks map[int]Task
}

func NewStorage() *Storage {
	return &Storage{
		tasks: make(map[int]Task),
	}
}

func (this *Storage) CreateTask(task Task) (*Task, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	task.TaskID = rand.Int()
	this.tasks[task.TaskID] = task

	return &task, nil
}
