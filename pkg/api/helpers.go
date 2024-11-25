package api

import (
	"encoding/json"
	"errors"
	"fmt"
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

func accessDenied(w http.ResponseWriter) error {
	return writeError(w, http.StatusUnauthorized, errors.New("access denied"))
}

func invalidId(w http.ResponseWriter) error {
	return writeError(w, http.StatusUnprocessableEntity, errors.New("invalid id"))
}

func conflict(w http.ResponseWriter) error {
	return writeError(w, http.StatusConflict, errors.New("conflict"))
}

func invalidRequest(w http.ResponseWriter) error {
	return writeError(w, http.StatusUnprocessableEntity, errors.New("invalid request"))
}

func notFound(w http.ResponseWriter) error {
	return writeError(w, http.StatusNotFound, errors.New("not found"))
}

func serverError(w http.ResponseWriter) error {
	return writeError(w, http.StatusInternalServerError, errors.New("something went wrong"))
}

func RouteHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			writeError(w, http.StatusNoContent, errors.New("no content"))
			return
		}

		if err := f(w, r); err != nil {
			writeError(w, http.StatusInternalServerError, err)
		}
	}
}

func writeJson(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&ApiResponse{Data: value})
}

func writeError(w http.ResponseWriter, status int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
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
