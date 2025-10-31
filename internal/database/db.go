package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Database connected!")
	applyMigrations()
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func applyMigrations() {
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			title TEXT NOT NULL,
			photo TEXT NOT NULL,
			likes INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER REFERENCES users(id),
			post_id INTEGER REFERENCES posts(id),
			PRIMARY KEY (user_id, post_id)
		);
	`

	_, err := DB.Exec(migrationSQL)
	if err != nil {
		log.Println("Migration error:", err)
	} else {
		log.Println("Migrations applied!")
	}
}