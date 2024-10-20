package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Implements Database Interface
type Store struct {
	db *sql.DB
}

func (s *Store) Connect() error {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", db_user, db_name, db_password)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Print("🗃️ connected to database")
	s.db = db
	return nil
}

func (s *Store) Init() {
	log.Println("Creating users table")
	query := `create table if not exists users (
	id serial primary key,
	first_name varchar(50),
	last_name varchar(50),
	email varchar(50),
	password varchar,
	created_at timestamp
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Println("Failed to create users table")
	}
	log.Println("Created users table")
}
