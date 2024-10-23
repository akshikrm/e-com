package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Registers user routes to the passed in router
func RegisterUserApi(r *http.ServeMux, db *db.PostgresStore) {
	userApi := &UserApi{
		UserService: user.NewUserService(db),
	}

	r.HandleFunc("GET /users", RouteHandler(userApi.GetAll))
	r.HandleFunc("POST /users", RouteHandler(userApi.Create))
	r.HandleFunc("GET /users/{id}", RouteHandler(userApi.GetOne))
	r.HandleFunc("DELETE /users/{id}", RouteHandler(userApi.Delete))
}

type UserApi struct {
	UserService types.UserService
}

func (u *UserApi) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := u.UserService.Get()
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, users)
}

func (u *UserApi) GetOne(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("invalid id")
	}

	foundUser, err := u.UserService.GetOne(id)

	if err != nil {
		if err == user.UserNotFound {
			return writeError(w, http.StatusNotFound, err)
		}
	}

	return writeJson(w, http.StatusOK, foundUser)
}

func (u *UserApi) Create(w http.ResponseWriter, r *http.Request) error {
	a := &types.User{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	err := u.UserService.Create(a)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusCreated, "created user")
}

func (u *UserApi) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("invalid id")
	}

	if err := u.UserService.Delete(id); err != nil {
		if err == user.UserNotFound {
			return writeError(w, http.StatusNotFound, err)
		}
		return err
	}

	return writeJson(w, http.StatusOK, "deleted successfully")
}

func parseId(id string) (int, error) {
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}

	return parsedId, nil

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
			writeError(w, http.StatusBadRequest, err)
		}
	}
}
