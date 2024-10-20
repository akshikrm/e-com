package main

import (
	db "akshidas/e-com/pkg/db"
	server "akshidas/e-com/pkg/server"
	"log"
	"os"
)

func main() {
	_, err := db.NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to connect to database %s", err)
		os.Exit(0)
	}

	server := &server.APIServer{Status: "Server is up and running", Port: ":5234"}
	server.Run()
	log.Printf("test")
}
