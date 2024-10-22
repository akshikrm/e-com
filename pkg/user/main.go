package user

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/types"
)

func NewUserService(db *db.Store) types.UserService {
	return &userService{DB: db.DB}
}
