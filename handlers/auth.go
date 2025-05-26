package handlers

import (
    "database/sql"
    "go-crud-app/models"
    "go-crud-app/views"
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "time"
)

// ShowLogin renders the login page
func ShowLogin() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        views.Login().Render(r.Context(), w)
    }
}

// ShowSignup renders the signup page
func ShowSignup() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        views.Signup().Render(r.Context(), w)
    }
}

// HandleLogin processes login form submission
func HandleLogin(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Get user from database
        var user models.User
        var hashedPassword string
        err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).Scan(
            &user.ID, &user.Name, &user.Email, &hashedPassword,
        )

        if err != nil {
            http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
            return
        }

        // Compare passwords
        err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
        if err != nil {
            http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
            return
        }

        // Set session cookie
        session, _ := generateSessionToken()
        http.SetCookie(w, &http.Cookie{
            Name:     "session_token",
            Value:    session,
            Path:     "/",
            Expires:  time.Now().Add(24 * time.Hour),
            HttpOnly: true,
        })

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

// HandleSignup processes signup form submission
func HandleSignup(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        name := r.FormValue("name")
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Check if email already exists
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if count > 0 {
            http.Redirect(w, r, "/signup?error=email_exists", http.StatusSeeOther)
            return
        }

        // Hash password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Insert new user
        _, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
            name, email, string(hashedPassword))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Redirect to login page
        http.Redirect(w, r, "/login?success=account_created", http.StatusSeeOther)
    }
}

// Helper function to generate a random session token
func generateSessionToken() (string, error) {
    // In a real app, use a proper session management library
    // This is just a simple example
    return "session_" + time.Now().Format("20060102150405"), nil
}