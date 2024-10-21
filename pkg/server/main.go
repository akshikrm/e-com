package server

import (
	"akshidas/e-com/pkg/types"
	"encoding/json"
	"log"
	"net/http"
)

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, http.StatusOK, "server is up and running")
	})

	router.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		a := &types.User{}
		if err := json.NewDecoder(r.Body).Decode(a); err != nil {
			writeJson(w, http.StatusBadRequest, err)
		}

		err := s.User.Create(a)
		if err != nil {
			writeJson(w, http.StatusBadRequest, "failed to add to new user")
		}

		writeJson(w, http.StatusOK, "created user")
	})

	log.Printf("🚀 Server started on port %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, router))
}

func writeJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
