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
	GetPasswordByEmail(email string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	Create(user types.CreateUserRequest) (*model.User, error)
	Update(id int, user types.UpdateUserRequest) error
	Delete(id int) error
}

type ProfileModeler interface {
	GetByUserId(int) (*model.Profile, error)
	Create(*types.NewProfileRequest) (int, error)
	UpdateProfileByUserID(int, *types.UpdateProfileRequest) error
	CheckIfUserExists(string) bool
}

type UserService struct {
	userModel    UserModeler
	profileModel ProfileModeler
}

func (u *UserService) Login(payload *types.LoginUserRequest) (string, error) {
	user, err := u.userModel.GetPasswordByEmail(payload.Email)
	if err != nil {
		return "", err
	}

	if err := utils.ValidateHash([]byte(user.Password), payload.Password); err != nil {
		log.Printf("invalid password for user %s", payload.Email)
		return "", utils.Unauthorized
	}

	token, err := utils.CreateJwt(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) Get() ([]*model.User, error) {
	return u.userModel.Get()
}

func (u *UserService) GetProfile(userId int) (*model.Profile, error) {
	return u.profileModel.GetByUserId(userId)
}

func (u *UserService) GetOne(id int) (*model.User, error) {
	user, err := u.userModel.GetOne(id)
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

	exists := u.profileModel.CheckIfUserExists(user.Email)
	if exists {
		return "", utils.Conflict
	}
	user.Password = hashedPassword
	savedUser, err := u.userModel.Create(user)
	if err != nil {
		return "", err
	}
	newUserProfile := types.NewProfileRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UserID:    savedUser.ID,
	}
	if _, err := u.profileModel.Create(&newUserProfile); err != nil {
		return "", err
	}
	fmt.Println(savedUser)
	token, err := utils.CreateJwt(savedUser.ID, savedUser.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserService) Update(id int, user *types.UpdateProfileRequest) (*model.Profile, error) {
	err := u.profileModel.UpdateProfileByUserID(id, user)
	if err != nil {
		return nil, err
	}
	return u.GetProfile(id)
}

func (u *UserService) Delete(id int) error {
	return u.userModel.Delete(id)
}

func NewUserService(userModel UserModeler, profileModel ProfileModeler) *UserService {
	return &UserService{userModel: userModel, profileModel: profileModel}
}
