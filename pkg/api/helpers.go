package api

import (
	"akshidas/e-com/pkg/utils"
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

func IsAuthenticated(f apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		authtoken := r.Header.Get("Authorization")
		token, err := utils.ValidateJWT(authtoken)
		if err != nil {
			accessDenied(w)
		}

		if !token.Valid {
			accessDenied(w)
		}
		return f(w, r)
	}

}

func accessDenied(w http.ResponseWriter) error {
	return writeError(w, http.StatusUnauthorized, errors.New("access denied"))

}

func RouteHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
