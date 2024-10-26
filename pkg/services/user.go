package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
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
	db UserModeler
}

func (u *UserService) Login(payload types.LoginUserRequest) (string, error) {
	user, err := u.db.GetUserByEmail(payload.Email)
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
	return u.db.Get()
}

func (u *UserService) GetOne(id int) (*model.User, error) {
	user, err := u.db.GetOne(id)
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
	userId, err := u.db.Create(user)
	if err != nil {
		return "", err
	}

	// userProfile := &types.NewProfileRequest{UserID: userId}
	// if err := u.profileService.Create(*userProfile); err != nil {
	// 	return "", err
	// }

	token, err := utils.CreateJwt(userId)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *UserService) Update(id int, user types.UpdateUserRequest) (*model.User, error) {
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
