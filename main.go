package main

import (
	"go-crud-app/database"
	"go-crud-app/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	r := chi.NewRouter()
	r.Get("/", handlers.ListTodos(db))
	r.Post("/create", handlers.CreateTodo(db))
	r.Post("/delete", handlers.DeleteTodo(db))
	r.Post("/toggle", handlers.ToggleTodo(db))

	http.ListenAndServe(":3000", r)
}
