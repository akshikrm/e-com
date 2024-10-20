package server

import (
	"akshidas/e-com/pkg/db"
	ecom "akshidas/e-com/pkg/e-com"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *APIServer) Run() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, s.Status)
	})

	e.POST("/users", Create(s.Store))

	e.Logger.Fatal(e.Start(s.Port))
}

func Create(s db.Store) echo.HandlerFunc {
	return func(c echo.Context) error {
		a := &ecom.User{}
		if err := json.NewDecoder(c.Request().Body).Decode(a); err != nil {
			return err
		}

		sqlQuery := `insert into udsers (first_name, last_name, password, email, created_at) values($1, $2, $3, $4, $5)`

		if _, err := s.DB.Query(sqlQuery,
			a.FirstName,
			a.LastName,
			a.Password,
			a.Email,
			time.Now().UTC(),
		); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())

		}
		return c.JSON(http.StatusCreated, "user created")
	}

}
