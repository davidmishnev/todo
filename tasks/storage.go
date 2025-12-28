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
	if task.Header == "" {
		return nil, ErrWrongArgument
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	task.TaskID = rand.Int()
	this.tasks[task.TaskID] = task
	return &task, nil
}

func (this *Storage) Delete(id int) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, exists := this.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(this.tasks, id)
	return nil
}

func (this *Storage) Update(id int, updated *Task) (*Task, error) {
	if updated.Header == "" {
		return nil, ErrWrongArgument
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()
	task, exists := this.tasks[id]
	if !exists {
		return nil, ErrWrongArgument
	}
	task.Header = updated.Header
	task.Description = updated.Description
	task.Status = updated.Status

	return &task, nil
}

func (this *Storage) GetAll() []*Task {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	result := make([]*Task, 0, len(this.tasks))
	for _, task := range this.tasks {
		result = append(result, &task)
	}

	return result
}

func (this *Storage) GetByID(id int) (*Task, error) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	task, exists := this.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}
