package api

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/user"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (u *UserApi) Update(w http.ResponseWriter, r *http.Request) error {
	a := &types.User{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	id, err := parseId(r.PathValue("id"))
	a.ID = id

	user, err := u.UserService.Update(a)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, user)
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
