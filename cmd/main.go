package main

import (
	db "akshidas/e-com/pkg/db"
	server "akshidas/e-com/pkg/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if _, err := db.NewPostgresStore(); err != nil {
		log.Fatalf("Failed to connect to database %s", err)
		os.Exit(0)
	}

	server := &server.APIServer{Status: "Server is up and running", Port: ":5234"}
	server.Run()
	log.Printf("test")
}
