package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type ApiResponse struct {
	Data any `json:"data"`
}

func RouteHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		if err := f(w, r); err != nil {
			writeError(w, http.StatusInternalServerError, err)
		}
	}
}

func writeJson(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&ApiResponse{Data: value})
}

func writeError(w http.ResponseWriter, status int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&ApiError{Error: err.Error()})
}

func parseId(id string) (int, error) {
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return parsedId, nil
}
