package api

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserServicer interface {
	Get() ([]*model.User, error)
	GetOne(id int) (*model.User, error)
	Create(user *model.User) (string, error)
	Update(user *model.User) (*model.User, error)
	Delete(id int) error
}

type UserApi struct {
	UserService UserServicer
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
		if err == utils.NotFound {
			return writeError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		}
		return err
	}

	return writeJson(w, http.StatusOK, foundUser)
}

func (u *UserApi) Create(w http.ResponseWriter, r *http.Request) error {
	a := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	token, err := u.UserService.Create(a)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusCreated, token)
}

func (u *UserApi) Update(w http.ResponseWriter, r *http.Request) error {
	a := &model.User{}
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
		if err == utils.NotFound {
			return writeError(w, http.StatusNotFound, err)
		}
		return err
	}

	return writeJson(w, http.StatusOK, "deleted successfully")
}

func NewUserApi(userService UserServicer) *UserApi {
	return &UserApi{UserService: userService}
}

// Registers user routes to the passed in router
