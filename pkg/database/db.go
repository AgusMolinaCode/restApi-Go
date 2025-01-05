package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = sql.Open("postgres", connStr)
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
			location_address TEXT NOT NULL,
			location_lng DOUBLE PRECISION NOT NULL,
			location_lat DOUBLE PRECISION NOT NULL,
			date_times JSONB NOT NULL,
			user_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			payment_link JSONB,
			tags TEXT[],
			transport_guide TEXT,
			schedule JSONB,
			exclusive_parking BOOLEAN DEFAULT FALSE,
			min_price DOUBLE PRECISION,
			rules JSONB,
			social_links JSONB,
			accessibility JSONB,
			delivery_method TEXT,
			main_image_url TEXT,
			additional_images JSONB,
			category TEXT NOT NULL
		);
	`

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			whatsapp TEXT NOT NULL UNIQUE,
			reset_token TEXT,
			reset_token_expiry TIMESTAMP,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	`

	createRegistrationsTable := `
		CREATE TABLE IF NOT EXISTS registrations (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			whatsapp TEXT,
			created_at TEXT NOT NULL,
			event_date TEXT,
			payment_link TEXT,
			FOREIGN KEY(event_id) REFERENCES events(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
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
