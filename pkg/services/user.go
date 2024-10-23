package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/utils"
)

type UserService struct {
	db *model.UserModel
}

func (u *UserService) Get() ([]*model.User, error) {
	return u.db.Get()
}

func (u *UserService) GetOne(id int) (*model.User, error) {
	return u.db.GetOne(id)
}

func (u *UserService) Create(user *model.User) (string, error) {
	hashedPassword, err := utils.HashPassword([]byte(user.Password))
	if err != nil {
		return "", err
	}

	user.Password = hashedPassword
	userId, err := u.db.Create(user)
	return utils.CreateJwt(userId)
}

func (u *UserService) Update(user *model.User) (*model.User, error) {
	err := u.db.Update(user)
	if err != nil {
		return nil, err
	}
	return u.GetOne(user.ID)
}

func (u *UserService) Delete(id int) error {
	return u.db.Delete(id)
}

func NewUserService(db *model.UserModel) *UserService {
	return &UserService{db: db}
}
