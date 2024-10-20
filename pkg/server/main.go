package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type APIServer struct {
	Status string
	Port   string
}

func (s *APIServer) Run() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, s.Status)
	})

	e.Logger.Fatal(e.Start(s.Port))
}
