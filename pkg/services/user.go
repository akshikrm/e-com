package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"fmt"
	"log"
)

type UserModeler interface {
	Get() ([]*model.User, error)
	GetOne(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	Create(user *types.CreateUserRequest) (int, error)
	Update(id int, user *types.UpdateUserRequest) error
	Delete(id int) error
}

type UserService struct {
	db UserModeler
}

func (u *UserService) Login(payolad *types.LoginUserRequest) (string, error) {

	user, err := u.db.GetUserByEmail(payolad.Email)

	if err != nil {
		return "", err
	}

	if err := utils.ValidateHash([]byte(user.Password), payolad.Password); err != nil {
		fmt.Println(err)
		log.Printf("invalid password for user %s", payolad.Email)
		return "", utils.Unauthorized
	}
	token, err := utils.CreateJwt(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}

func (u *UserService) Get() ([]*model.User, error) {
	return u.db.Get()
}

func (u *UserService) GetOne(id int) (*model.User, error) {
	return u.db.GetOne(id)
}

func (u *UserService) Create(user *types.CreateUserRequest) (string, error) {
	hashedPassword, err := utils.HashPassword([]byte(user.Password))
	if err != nil {
		return "", err
	}

	user.Password = hashedPassword
	userId, err := u.db.Create(user)
	return utils.CreateJwt(userId)
}

func (u *UserService) Update(id int, user *types.UpdateUserRequest) (*model.User, error) {
	err := u.db.Update(id, user)
	if err != nil {
		return nil, err
	}
	return u.GetOne(id)
}

func (u *UserService) Delete(id int) error {
	return u.db.Delete(id)
}

func NewUserService(db UserModeler) *UserService {
	return &UserService{db: db}
}
