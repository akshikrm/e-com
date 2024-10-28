package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
)

type UserModeler interface {
	Get() ([]*model.User, error)
	GetOne(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	Create(user types.CreateUserRequest) (int, error)
	Update(id int, user types.UpdateUserRequest) error
	Delete(id int) error
}

type UserService struct {
	db    *sql.DB
	model UserModeler
}

func (u *UserService) Login(payload types.LoginUserRequest) (string, error) {
	user, err := u.model.GetUserByEmail(payload.Email)
	if err != nil {
		return "", err
	}

	if err := utils.ValidateHash([]byte(user.Password), payload.Password); err != nil {
		log.Printf("invalid password for user %s", payload.Email)
		return "", utils.Unauthorized
	}

	token, err := utils.CreateJwt(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) Get() ([]*model.User, error) {
	return u.model.Get()
}

func (u *UserService) GetOne(id int) (*model.User, error) {
	user, err := u.model.GetOne(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) Create(user types.CreateUserRequest) (string, error) {
	hashedPassword, err := utils.HashPassword([]byte(user.Password))
	if err != nil {
		return "", err
	}

	user.Password = hashedPassword
	userId, err := u.model.Create(user)
	if err != nil {
		return "", err
	}

	userProfileService := NewProfileService(u.db)
	newUserProfile := &types.NewProfileRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UserID:    userId,
	}

	if err := userProfileService.Create(newUserProfile); err != nil {
		return "", err
	}

	token, err := utils.CreateJwt(userId)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *UserService) Update(id int, user types.UpdateUserRequest) (*model.User, error) {
	err := u.model.Update(id, user)
	if err != nil {
		return nil, err
	}
	return u.GetOne(id)
}

func (u *UserService) Delete(id int) error {
	return u.model.Delete(id)
}

func NewUserService(database *sql.DB) *UserService {
	userModel := model.NewUserModel(database)
	return &UserService{model: userModel, db: database}
}
