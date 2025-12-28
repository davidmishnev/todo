package tasks

import (
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
	storage := NewStorage()

	task := Task{Header: "Test", Description: "Description"}
	_, err := storage.CreateTask(task)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateEmptyTitle(t *testing.T) {
	storage := NewStorage()

	task := Task{Header: "", Description: "Description"}
	_, err := storage.CreateTask(task)

	if !errors.Is(err, ErrWrongArgument) {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestGetByID(t *testing.T) {
	storage := NewStorage()
	task := Task{Header: "Test"}
	created, _ := storage.CreateTask(task)

	found, err := storage.GetByID(created.TaskID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Header != "Test" {
		t.Errorf("expected title 'Test', got %storage", found.Header)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	storage := NewStorage()

	_, err := storage.GetByID(999)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	storage := NewStorage()
	task := Task{Header: "Original"}
	created, _ := storage.CreateTask(task)

	updated, err := storage.Update(created.TaskID, &Task{Header: "Updated", Status: 0})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Header != "Updated" {
		t.Errorf("expected title 'Updated', got %storage", updated.Header)
	}
}

func TestDelete(t *testing.T) {
	storage := NewStorage()
	task := Task{Header: "Test"}
	created, _ := storage.CreateTask(task)

	err := storage.Delete(created.TaskID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = storage.GetByID(created.TaskID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestCreateWithStatus(t *testing.T) {
	s := NewStorage()

	tests := []struct {
		name   string
		status TaskStatus
	}{
		{"assigned status", Assigned},
		{"in progress status", InProgress},
		{"completed status", Completed},
		{"dropped status", Dropped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{Header: "Test", Description: "Description", Status: tt.status}
			created, err := s.CreateTask(task)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if created.Status != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, created.Status)
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	s := NewStorage()
	task := Task{Header: "Test", Status: Assigned}
	created, _ := s.CreateTask(task)

	tests := []struct {
		name      string
		newStatus TaskStatus
	}{
		{"change to in progress", InProgress},
		{"change to completed", Completed},
		{"change to dropped", Dropped},
		{"change back to assigned", Assigned},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := s.Update(created.TaskID, &Task{Header: "Test", Status: tt.newStatus})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if updated.Status != tt.newStatus {
				t.Errorf("expected status %d, got %d", tt.newStatus, updated.Status)
			}
		})
	}
}

func TestGetAllWithDifferentStatuses(t *testing.T) {
	s := NewStorage()

	tasks := []Task{
		{Header: "Task 1", Status: Assigned},
		{Header: "Task 2", Status: InProgress},
		{Header: "Task 3", Status: Completed},
		{Header: "Task 4", Status: Dropped},
	}

	for _, task := range tasks {
		_, err := s.CreateTask(task)
		if err != nil {
			t.Fatalf("failed to create task: %v", err)
		}
	}

	all := s.GetAll()
	if len(all) != 4 {
		t.Errorf("expected 4 tasks, got %d", len(all))
	}

	statusCount := make(map[TaskStatus]int)
	for _, task := range all {
		statusCount[task.Status]++
	}

	if statusCount[Assigned] != 1 || statusCount[InProgress] != 1 ||
		statusCount[Completed] != 1 || statusCount[Dropped] != 1 {
		t.Errorf("expected 1 task of each status, got %v", statusCount)
	}
}

func TestUpdateStatusOtherFields(t *testing.T) {
	s := NewStorage()
	task := Task{Header: "Original", Description: "Original Description", Status: Assigned}
	created, _ := s.CreateTask(task)

	updated, err := s.Update(created.TaskID, &Task{
		Header:      "Original",
		Description: "Original Description",
		Status:      Completed,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Header != "Original" {
		t.Errorf("expected header 'Original', got %s", updated.Header)
	}
	if updated.Description != "Original Description" {
		t.Errorf("expected description 'Original Description', got %s", updated.Description)
	}
	if updated.Status != Completed {
		t.Errorf("expected status Completed, got %d", updated.Status)
	}
}
