package api

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserServicer interface {
	Get() ([]*model.User, error)
	GetOne(id int) (*model.User, error)
	Login(types.LoginUserRequest) (string, error)
	Create(user types.CreateUserRequest) (string, error)
	Update(id int, user types.UpdateUserRequest) (*model.User, error)
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

func (u *UserApi) Login(w http.ResponseWriter, r *http.Request) error {
	a := &types.LoginUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	token, err := u.UserService.Login(*a)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, token)
}

func (u *UserApi) Create(w http.ResponseWriter, r *http.Request) error {
	a := &types.CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	token, err := u.UserService.Create(*a)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusCreated, token)
}

func (u *UserApi) Update(w http.ResponseWriter, r *http.Request) error {
	a := &types.UpdateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	id, err := parseId(r.PathValue("id"))

	user, err := u.UserService.Update(id, *a)
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
