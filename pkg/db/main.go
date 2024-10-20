package db

import (
	"log"
)

type Database interface {
	Connect() error
	Init() error
}

func Connect(d Database) {
	err := d.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database %s", err)
	}
	d.Init()

}
