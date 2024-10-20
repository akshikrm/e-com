package server

import (
	db "akshidas/e-com/pkg/db"
)

type APIServer struct {
	Status string
	Port   string
	Store  db.Store
}
