package main

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database := &db.PostgresStore{}
	db.Connect(database)

	server := &server.APIServer{
		Status: "Server is up and running",
		Port:   ":5234",
		Store:  database,
	}
	server.Run()
}
