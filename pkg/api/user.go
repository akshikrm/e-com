package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/storage"
	"akshidas/e-com/pkg/types"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type UserServicer interface {
	Get() ([]*types.User, error)
	GetProfile(int) (*types.Profile, error)
	GetOne(int) (*types.User, error)
	Login(*types.LoginUserRequest) (string, error)
	Create(types.CreateUserRequest) (string, error)
	Update(int, *types.UpdateProfileRequest) (*types.Profile, error)
	Delete(int) error
}

type UserApi struct {
	UserService UserServicer
}

type UserProfile struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	Profile   *types.Profile `json:"profile"`
}

func (u *UserApi) GetProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := ctx.Value("userID")
	userProfile, err := u.UserService.GetProfile(id.(int))
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, userProfile)
}

func (u *UserApi) UpdateProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := ctx.Value("userID")
	a := types.UpdateProfileRequest{}
	if err := DecodeBody(r.Body, &a); err != nil {
		return err
	}
	user, err := u.UserService.Update(id.(int), &a)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, user)
}

func (u *UserApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := u.UserService.Get()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, users)
}
func (u *UserApi) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	foundUser, err := u.UserService.GetOne(id)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, foundUser)
}

func (u *UserApi) Login(w http.ResponseWriter, r *http.Request) error {
	a := types.LoginUserRequest{}
	if err := DecodeBody(r.Body, &a); err != nil {
		return err
	}
	token, err := u.UserService.Login(&a)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, token)
}

func (u *UserApi) Create(w http.ResponseWriter, r *http.Request) error {
	a := &types.CreateUserRequest{}
	if err := DecodeBody(r.Body, &a); err != nil {
		return err
	}
	token, err := u.UserService.Create(*a)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, token)
}

func (u *UserApi) Update(w http.ResponseWriter, r *http.Request) error {
	a := types.UpdateProfileRequest{}
	if err := DecodeBody(r.Body, &a); err != nil {
		return err
	}
	id, err := parseId(r.PathValue("id"))
	user, err := u.UserService.Update(id, &a)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, user)
}

func (u *UserApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	log.Println("deleting")
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	if err := u.UserService.Delete(id); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, "deleted successfully")
}

func NewUserApi(store *db.Storage) *UserApi {
	userModel := storage.NewUserStorage(store.DB)
	profileModel := storage.NewProfileStorage(store.DB)
	userService := services.NewUserService(userModel, profileModel)
	return &UserApi{UserService: userService}
}
