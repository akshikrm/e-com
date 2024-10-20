package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type APIServer struct {
	Status string
	Port   string
}

type User struct {
	FistName string
	LastName string
	Email    string
	Password string
}

func (s *APIServer) Run() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, s.Status)
	})

	e.GET("/users", GetAllUsers)

	e.Logger.Fatal(e.Start(s.Port))
}

func GetAllUsers(c echo.Context) error {
	users := &[]User{{FistName: "Akshay", LastName: "Krishna", Email: "akshay@bpract.com", Password: "root"}, {FistName: "Akshay", LastName: "Krishna", Email: "akshay@bpract.com", Password: "root"}}

	return c.JSON(http.StatusCreated, users)
}
