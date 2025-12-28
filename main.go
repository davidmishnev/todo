package main

import (
	"log"
	"net/http"
	"todo/tasks"
)

func main() {
	storage := tasks.NewStorage()
	logger := log.Default()
	server := NewServer(storage, logger)

	http.HandleFunc("/todos", server.loggingMiddleware(server.handleTodos))
	http.HandleFunc("/todos/", server.loggingMiddleware(server.handleTodoByID))

	logger.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal(err)
	}
}
