package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/user"
	"encoding/json"
	"net/http"
)

// Registers user routes to the passed in router
func RegisterUserApi(r *http.ServeMux, db *db.PostgresStore) {
	userApi := &UserApi{
		UserService: user.NewUserService(db),
	}
	r.HandleFunc("GET /users", userApi.GetAll)
	r.HandleFunc("POST /users", userApi.Create)
}

type UserApi struct {
	UserService types.UserService
}

func (u *UserApi) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.UserService.Get()
	if err != nil {
		writeJson(w, http.StatusBadRequest, "failed to get data")
		return
	}
	writeJson(w, http.StatusOK, users)
}

func (u *UserApi) Create(w http.ResponseWriter, r *http.Request) {
	a := &types.User{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		writeJson(w, http.StatusBadRequest, err)
	}

	err := u.UserService.Create(a)
	if err != nil {
		writeJson(w, http.StatusBadRequest, "failed to add to new user")
	}

	writeJson(w, http.StatusOK, "created user")
}

func writeJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
