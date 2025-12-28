package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo/tasks"
)

const SecToTimeout = 5

type Server struct {
	storage *tasks.Storage
	logger  *log.Logger
}

func NewServer(storage *tasks.Storage, logger *log.Logger) *Server {
	return &Server{
		storage: storage,
		logger:  logger,
	}
}

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.logger.Printf("%s %s started", r.Method, r.URL.Path)
		next(w, r)
		s.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func (s *Server) handleTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), SecToTimeout*time.Second)
	defer cancel()

	r = r.WithContext(ctx)

	switch r.Method {
	case http.MethodPost:
		s.createTodo(w, r)
	case http.MethodGet:
		s.getAllTodos(w)
	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTodoByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), SecToTimeout*time.Second)
	defer cancel()

	r = r.WithContext(ctx)

	id, err := s.extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getTodoByID(w, id)
	case http.MethodPut:
		s.updateTodo(w, r, id)
	case http.MethodDelete:
		s.deleteTodo(w, id)
	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) extractID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, http.ErrNoLocation
	}
	return strconv.Atoi(parts[1])
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	created, err := s.storage.CreateTask(task)
	if err != nil {
		switch {
		case errors.Is(err, tasks.ErrWrongArgument):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(created)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getAllTodos(w http.ResponseWriter) {
	tasksGot := s.storage.GetAll()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasksGot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getTodoByID(w http.ResponseWriter, id int) {
	task, err := s.storage.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, tasks.ErrTaskNotFound):
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request, id int) {
	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := s.storage.Update(id, &task)
	if err != nil {
		switch {
		case errors.Is(err, tasks.ErrWrongArgument):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, tasks.ErrTaskNotFound):
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updated)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteTodo(w http.ResponseWriter, id int) {
	err := s.storage.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, tasks.ErrTaskNotFound):
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
