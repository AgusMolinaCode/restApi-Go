package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	DB.SetMaxOpenConns(8)
	DB.SetMaxIdleConns(6)

	createTables()
}

func createTables() {
	createEventsTable := `
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			date_time TEXT NOT NULL,
			user_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			payment_link TEXT
		)
	`

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			reset_token TEXT,
			reset_token_expiry DATETIME
		)
	`

	createRegistrationsTable := `
		CREATE TABLE IF NOT EXISTS registrations (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY(event_id) REFERENCES events(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`

	_, err := DB.Exec(createEventsTable)
	if err != nil {
		log.Fatalf("Error creating events table: %v", err)
	}

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		log.Fatalf("Error creating registrations table: %v", err)
	}
}
