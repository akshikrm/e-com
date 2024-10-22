package server

import (
	"akshidas/e-com/pkg/api"
	"akshidas/e-com/pkg/db"
	"log"
	"net/http"
)

type APIServer struct {
	Status string
	Port   string
	Store  *db.PostgresStore
}

// Create a new server and registers routes to it
func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is up and running"))
	})

	api.RegisterUserApi(router, s.Store)
	log.Printf("🚀 Server started on port %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, router))
}
