package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Storage struct {
	DB       *sql.DB
	user     string
	name     string
	password string
	host     string
	port     string
}

func Connect(s *Storage) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		s.host, s.port, s.user, s.name, s.password)
	db, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	s.DB = db
	log.Println("🗃️ connected to database")
}

func NewStorage() *Storage {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return &Storage{
		user:     db_user,
		name:     db_name,
		password: db_password,
		host:     db_host,
		port:     db_port,
	}
}
