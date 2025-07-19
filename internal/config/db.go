package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Set reasonable connection pool limits for high concurrency
	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Check connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL!")

	// Create tables if not exist (basic migration)
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

// CreateTables runs basic SQL migration
func createTables() error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	);`

	postTable := `
        CREATE TABLE IF NOT EXISTS posts (
                id SERIAL PRIMARY KEY,
                user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                content TEXT NOT NULL,
                likes INT DEFAULT 0,
                created_at TIMESTAMP DEFAULT now()
        );`

	likeTable := `
        CREATE TABLE IF NOT EXISTS likes (
                id SERIAL PRIMARY KEY,
                user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                created_at TIMESTAMP DEFAULT now(),
                UNIQUE(user_id, post_id)
        );`

	followTable := `
	CREATE TABLE IF NOT EXISTS follows (
		id SERIAL PRIMARY KEY,
		follower_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		followed_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT now()
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		return err
	}
	_, err = DB.Exec(postTable)
	if err != nil {
		return err
	}
	_, err = DB.Exec(followTable)
	if err != nil {
		return err
	}

	_, err = DB.Exec(likeTable)
	if err != nil {
		return err
	}

	return nil
}
