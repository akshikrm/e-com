package db

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

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

	seed := flag.Bool("initdb", false, "initialze db if true")

	flag.Parse()
	if *seed {
		s.Init()
	}

	return nil
}

func (s *PostgresStore) Init() {
	CreateUserTable(s.DB)
	CreateProfileTable(s.DB)
	log.Println("successfully created all tables")
	os.Exit(0)
}

func CreateUserTable(db *sql.DB) {
	log.Println("Creating users table")
	query := `CREATE TABLE IF NOT EXISTS users (
	id serial primary key,
	password varchar,
	created_at timestamp
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create users table %s", err)
		return
	}
	log.Println("Created users table")
}

func CreateProfileTable(db *sql.DB) {
	log.Println("Creating profiles table")
	query := `CREATE TABLE IF NOT EXISTS profiles (
	id serial primary key,
	user_id int,
	first_name varchar(50) DEFAULT '' NOT NULL,
	last_name varchar(50) DEFAULT '' NOT NULL,
	email varchar(50) DEFAULT '' NOT NULL,
	pincode varchar(10) DEFAULT '' NOT NULL,
	address_one varchar(100) DEFAULT '' NOT NULL,
	address_two varchar(100) DEFAULT '' NOT NULL,
	phone_number varchar(15) DEFAULT '' NOT NULL,
	created_at timestamp,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create profiles table %s", err)
		return
	}
	log.Println("Created profiles table")
}

func (s *PostgresStore) getConnectionString() string {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", db_host, db_port, db_user, db_name, db_password)

}
