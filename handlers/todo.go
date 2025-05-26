package handlers

import (
	"database/sql"
	"fmt"
	"go-crud-app/models"
	"go-crud-app/views"
	"net/http"
	"strconv"
)

// ListTodos renders todos with pagination (5 items per page)
func ListTodos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=== ListTodos called ===")

		// Parse page parameter
		pageStr := r.URL.Query().Get("page")
		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		const itemsPerPage = 5
		offset := (page - 1) * itemsPerPage

		// Get total count
		var totalCount int
		err := db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&totalCount)
		if err != nil {
			fmt.Printf("ERROR: Failed to count todos: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate pagination data
		totalPages := (totalCount + itemsPerPage - 1) / itemsPerPage
		if totalPages == 0 {
			totalPages = 1
		}
		if page > totalPages {
			page = totalPages
		}

		paginationData := models.PaginationData{
			CurrentPage:  page,
			TotalPages:   totalPages,
			HasPrevious:  page > 1,
			HasNext:      page < totalPages,
			PreviousPage: page - 1,
			NextPage:     page + 1,
		}

		// Get paginated todos
		rows, err := db.Query("SELECT id, task, done FROM todos ORDER BY id DESC LIMIT ? OFFSET ?", itemsPerPage, offset)
		if err != nil {
			fmt.Printf("ERROR: Failed to query todos: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []models.Todo
		for rows.Next() {
			var t models.Todo
			err := rows.Scan(&t.ID, &t.Task, &t.Done)
			if err != nil {
				fmt.Printf("ERROR: Failed to scan todo: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			todos = append(todos, t)
		}

		if err = rows.Err(); err != nil {
			fmt.Printf("ERROR: Row iteration error: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Found %d todos (page %d of %d, total: %d)\n", len(todos), page, totalPages, totalCount)
		views.Todos(todos, paginationData).Render(r.Context(), w)
	}
}

// CreateTodo adds a new todo
func CreateTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=== CreateTodo called ===")

		// Ensure we're using POST method
		if r.Method != http.MethodPost {
			fmt.Println("ERROR: Only POST method is allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("ERROR: Failed to parse form: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the task value
		task := r.FormValue("task")
		fmt.Printf("Task received: '%s'\n", task)

		// Dump all form values for debugging
		fmt.Println("All form values:")
		for key, values := range r.Form {
			fmt.Printf("  %s: %v\n", key, values)
		}

		// Validate task
		if task == "" {
			fmt.Println("WARNING: Empty task received, nothing inserted")
			http.Redirect(w, r, "/?error=empty_task", http.StatusSeeOther)
			return
		}

		// Insert the task
		result, err := db.Exec("INSERT INTO todos (task, done) VALUES (?, ?)", task, false)
		if err != nil {
			fmt.Printf("ERROR: Failed to insert todo: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		lastInsertID, _ := result.LastInsertId()
		fmt.Printf("Todo created successfully. ID: %d, Rows affected: %d\n", lastInsertID, rowsAffected)

		fmt.Println("=== Redirecting back to / ===")
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
