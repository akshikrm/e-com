package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Print("🗃️ connected to database")
	return &PostgresStore{db: db}, nil
}
