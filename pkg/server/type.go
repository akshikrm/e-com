package server

import (
	"akshidas/e-com/pkg/types"
)

type APIServer struct {
	Status string
	Port   string
	User   types.UserService
}
