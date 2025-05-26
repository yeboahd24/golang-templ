package handlers

import (
	"database/sql"
	"fmt"
	"go-crud-app/models"
	"go-crud-app/views"
	"net/http"
	"strconv"
)

// ListTodos renders all todos
func ListTodos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, task, done FROM todos")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []models.Todo
		for rows.Next() {
			var t models.Todo
			err := rows.Scan(&t.ID, &t.Task, &t.Done)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			todos = append(todos, t)
		}

		views.Todos(todos).Render(r.Context(), w)
	}
}

// CreateTodo adds a new todo
func CreateTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		task := r.FormValue("task")

		if task != "" {
			_, err := db.Exec("INSERT INTO todos (task, done) VALUES (?, ?)", task, false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// DeleteTodo removes a todo by ID
func DeleteTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM todos WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ToggleTodo toggles the done state of a todo
func ToggleTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=== ToggleTodo called ===")

		r.ParseForm()
		idStr := r.FormValue("id")
		fmt.Printf("ID received: '%s'\n", idStr)

		if idStr == "" {
			fmt.Println("ERROR: No ID provided")
			http.Error(w, "No ID provided", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("ERROR: Invalid ID: %v\n", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		fmt.Printf("Parsed ID: %d\n", id)

		// Get current state
		var currentDone bool
		err = db.QueryRow("SELECT done FROM todos WHERE id = ?", id).Scan(&currentDone)
		if err != nil {
			fmt.Printf("ERROR: Failed to get current state: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Current done state: %v\n", currentDone)

		// Calculate new state
		newDone := !currentDone
		fmt.Printf("New done state will be: %v\n", newDone)

		// Update with explicit value instead of NOT
		result, err := db.Exec("UPDATE todos SET done = ? WHERE id = ?", newDone, id)
		if err != nil {
			fmt.Printf("ERROR: Failed to update: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("Rows affected: %d\n", rowsAffected)

		if rowsAffected == 0 {
			fmt.Printf("WARNING: No rows were updated - ID %d might not exist\n", id)
		}

		// Verify the change
		var verifyDone bool
		err = db.QueryRow("SELECT done FROM todos WHERE id = ?", id).Scan(&verifyDone)
		if err != nil {
			fmt.Printf("ERROR: Failed to verify update: %v\n", err)
		} else {
			fmt.Printf("Verified done state after update: %v\n", verifyDone)
		}

		fmt.Println("=== Redirecting back to / ===")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
