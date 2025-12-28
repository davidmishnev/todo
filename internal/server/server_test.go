package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"todo/internal/storage"
)

func setupServer() *Server {
	st := storage.NewStorage()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	return NewServer(st, logger)
}

func TestCreateTodo(t *testing.T) {
	server := setupServer()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
	}{
		{
			name:           "successful creation",
			payload:        `{"Header":"Test Task","Description":"Test Description","Status":0}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "empty header validation error",
			payload:        `{"Header":"","Description":"Test Description","Status":0}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			payload:        `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(tt.payload))
			w := httptest.NewRecorder()

			server.HandleTodos(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetAllTodos(t *testing.T) {
	server := setupServer()

	task := storage.Task{Header: "Test Task", Description: "Test Description", Status: storage.Assigned}
	_, _ = server.storage.CreateTask(task)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	server.HandleTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []storage.Task
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 task, got %d", len(result))
	}
}

func TestGetTodoByID(t *testing.T) {
	server := setupServer()

	task := storage.Task{Header: "Test Task", Description: "Test Description", Status: storage.Assigned}
	created, _ := server.storage.CreateTask(task)

	tests := []struct {
		name           string
		id             int
		expectedStatus int
	}{
		{
			name:           "successful get",
			id:             created.TaskID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "task not found",
			id:             999999,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/todos/"+strconv.Itoa(tt.id), nil)
			w := httptest.NewRecorder()

			server.HandleTodoByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	server := setupServer()

	task := storage.Task{Header: "Original Task", Description: "Original Description", Status: storage.Assigned}
	created, _ := server.storage.CreateTask(task)

	tests := []struct {
		name           string
		id             int
		payload        string
		expectedStatus int
	}{
		{
			name:           "successful update",
			id:             created.TaskID,
			payload:        `{"Header":"Updated Task","Description":"Updated Description","Status":2}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty header validation error",
			id:             created.TaskID,
			payload:        `{"Header":"","Description":"Updated Description","Status":2}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "task not found",
			id:             999999,
			payload:        `{"Header":"Updated Task","Description":"Updated Description","Status":2}`,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/todos/"+strconv.Itoa(tt.id), bytes.NewBufferString(tt.payload))
			w := httptest.NewRecorder()

			server.HandleTodoByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUpdateTodoStatus(t *testing.T) {
	server := setupServer()

	task := storage.Task{Header: "Test Task", Description: "Test", Status: storage.Assigned}
	created, _ := server.storage.CreateTask(task)

	tests := []struct {
		name           string
		status         storage.TaskStatus
		expectedStatus int
	}{
		{"update to in progress", storage.InProgress, http.StatusOK},
		{"update to completed", storage.Completed, http.StatusOK},
		{"update to dropped", storage.Dropped, http.StatusOK},
		{"update back to assigned", storage.Assigned, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := `{"Header":"Test Task","Description":"Test","Status":` + strconv.Itoa(int(tt.status)) + `}`
			req := httptest.NewRequest(http.MethodPut, "/todos/"+strconv.Itoa(created.TaskID), bytes.NewBufferString(payload))
			w := httptest.NewRecorder()

			server.HandleTodoByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var result storage.Task
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if result.Status != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, result.Status)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	server := setupServer()

	task := storage.Task{Header: "Test Task", Description: "Test Description", Status: storage.Assigned}
	created, _ := server.storage.CreateTask(task)

	tests := []struct {
		name           string
		id             int
		expectedStatus int
	}{
		{
			name:           "successful delete",
			id:             created.TaskID,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "task not found",
			id:             999999,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(tt.id), nil)
			w := httptest.NewRecorder()

			server.HandleTodoByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := setupServer()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"PATCH on /todos", http.MethodPatch, "/todos"},
		{"PATCH on /todos/1", http.MethodPatch, "/todos/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			if tt.path == "/todos" {
				server.HandleTodos(w, req)
			} else {
				server.HandleTodoByID(w, req)
			}

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
			}
		})
	}
}

func TestInvalidID(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/todos/invalid", nil)
	w := httptest.NewRecorder()

	server.HandleTodoByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetAllTodosWithDifferentStatuses(t *testing.T) {
	server := setupServer()

	testTasks := []storage.Task{
		{Header: "Task 1", Description: "Desc 1", Status: storage.Assigned},
		{Header: "Task 2", Description: "Desc 2", Status: storage.InProgress},
		{Header: "Task 3", Description: "Desc 3", Status: storage.Completed},
		{Header: "Task 4", Description: "Desc 4", Status: storage.Dropped},
	}

	for _, task := range testTasks {
		_, _ = server.storage.CreateTask(task)
	}

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	server.HandleTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []storage.Task
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("expected 4 tasks, got %d", len(result))
	}

	statusCount := make(map[storage.TaskStatus]int)
	for _, task := range result {
		statusCount[task.Status]++
	}

	if statusCount[storage.Assigned] != 1 || statusCount[storage.InProgress] != 1 ||
		statusCount[storage.Completed] != 1 || statusCount[storage.Dropped] != 1 {
		t.Errorf("expected 1 task of each status, got %v", statusCount)
	}
}

func TestCreateTodoWithDifferentStatuses(t *testing.T) {
	server := setupServer()

	tests := []struct {
		name           string
		status         storage.TaskStatus
		expectedStatus int
	}{
		{"create with assigned status", storage.Assigned, http.StatusCreated},
		{"create with in progress status", storage.InProgress, http.StatusCreated},
		{"create with completed status", storage.Completed, http.StatusCreated},
		{"create with dropped status", storage.Dropped, http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := `{"Header":"Test","Description":"Desc","Status":` + strconv.Itoa(int(tt.status)) + `}`
			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(payload))
			w := httptest.NewRecorder()

			server.HandleTodos(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var result storage.Task
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if result.Status != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, result.Status)
			}
		})
	}
}
