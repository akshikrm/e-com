package api

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserServicer interface {
	Get() ([]*model.User, error)
	GetProfile(int) (*model.Profile, error)
	GetOne(int) (*model.User, error)
	Login(types.LoginUserRequest) (string, error)
	Create(types.CreateUserRequest) (string, error)
	Update(int, types.UpdateUserRequest) (*model.User, error)
	Delete(int) error
}

type ProfileServicer interface {
	GetByUserId(int) (*model.Profile, error)
	Create(*types.NewProfileRequest) error
	Update(int, *types.UpdateProfileRequest) (*model.Profile, error)
}

type UserApi struct {
	UserService    UserServicer
	profileService ProfileServicer
}

type UserProfile struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	Profile   *model.Profile `json:"profile"`
}

func (u *UserApi) GetProfile(id int, w http.ResponseWriter, r *http.Request) error {
	// user, err := u.UserService.GetOne(id)
	// if err != nil {
	// 	return err
	// }
	//
	// profile, err := u.profileService.GetByUserId(id)
	// if err != nil {
	// 	return err
	// }
	//
	// userProfile := &UserProfile{
	// 	CreatedAt: user.CreatedAt,
	// 	Profile:   profile,
	// }

	userProfile, err := u.UserService.GetProfile(id)
	if err != nil {
		if err == utils.NotFound {
			return writeJson(w, http.StatusNotFound, "not found")
		}
		return err
	}
	return writeJson(w, http.StatusOK, userProfile)

}

func (u *UserApi) UpdateProfile(userId int, w http.ResponseWriter, r *http.Request) error {
	a := &types.UpdateProfileRequest{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		return err
	}

	user, err := u.profileService.Update(userId, a)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, user)
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
		if err == io.EOF {
			return errors.New("invalid request")
		}
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
		if err == io.EOF {
			return errors.New("invalid request")
		}
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
