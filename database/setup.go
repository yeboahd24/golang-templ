package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTables(db *sql.DB) {
	log.Println("Tables created")

	// Create tables
	db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		done BOOLEAN NOT NULL
	)`)
}

func SeedData(db *sql.DB) {
	log.Println("Seed data")
	// Seed data
	db.Exec(`INSERT INTO todos (task, done) VALUES ('Buy groceries', 0)`)
	db.Exec(`INSERT INTO todos (task, done) VALUES ('Clean the house', 0)`)
	db.Exec(`INSERT INTO todos (task, done) VALUES ('Call mom', 0)`)
}

func InitDB() *sql.DB {
	log.Println("Init DB")
	// Open database
	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	CreateTables(db)

	// Seed data
	SeedData(db)

	return db
}
