package user

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/types"
)

func NewUserService(db *db.PostgresStore) types.UserService {
	return &userService{DB: db.DB}
}
