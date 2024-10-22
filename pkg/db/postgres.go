package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// Implements Database Interface
type PostgresStore struct {
	DB *sql.DB
}

func (s *PostgresStore) Connect() error {
	db, err := sql.Open("postgres", s.getConnectionString())

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Print("🗃️ connected to database")
	s.DB = db
	return nil
}

func (s *PostgresStore) Init() {
	log.Println("Creating users table")

	query := `create table if not exists users (
	id serial primary key,
	first_name varchar(50),
	last_name varchar(50),
	email varchar(50),
	password varchar,
	created_at timestamp
	)`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create users table")
	}

	log.Println("Created users table")
}

func (s *PostgresStore) getConnectionString() string {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", db_host, db_port, db_user, db_name, db_password)

}
