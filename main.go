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

	// Todo routes
	r.Get("/", handlers.ListTodos(db))
	r.Post("/create", handlers.CreateTodo(db))
	r.Post("/delete", handlers.DeleteTodo(db))
	r.Post("/toggle", handlers.ToggleTodo(db))

	// Auth routes
	r.Get("/login", handlers.ShowLogin())
	r.Post("/login", handlers.HandleLogin(db))
	r.Get("/signup", handlers.ShowSignup())
	r.Post("/signup", handlers.HandleSignup(db))

	http.ListenAndServe(":3000", r)
}
