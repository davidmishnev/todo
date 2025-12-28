package tasks

import (
	"sync"
)

type Storage struct {
	counter int
	mutex   sync.RWMutex
	tasks   map[int]Task
}

func NewStorage() *Storage {
	return &Storage{
		counter: 0,
		tasks:   make(map[int]Task),
	}
}

func (s *Storage) CreateTask(task Task) (*Task, error) {
	if task.Header == "" {
		return nil, ErrWrongArgument
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	task.TaskID = s.counter
	s.counter++
	s.tasks[task.TaskID] = task

	return &task, nil
}

func (s *Storage) Delete(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

func (s *Storage) Update(id int, updated *Task) (*Task, error) {
	if updated.Header == "" {
		return nil, ErrWrongArgument
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	task.Header = updated.Header
	task.Description = updated.Description
	task.Status = updated.Status
	s.tasks[id] = task

	return &task, nil
}

func (s *Storage) GetAll() []Task {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}

	return result
}

func (s *Storage) GetByID(id int) (Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}
