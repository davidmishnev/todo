package main

import (
	"log"
	"net/http"
	"todo/internal/server"
	"todo/internal/storage"
)

func main() {
	st := storage.NewStorage()
	logger := log.Default()
	srv := server.NewServer(st, logger)

	http.HandleFunc("/todos", srv.LoggingMiddleware(srv.HandleTodos))
	http.HandleFunc("/todos/", srv.LoggingMiddleware(srv.HandleTodoByID))

	logger.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal(err)
	}
}
