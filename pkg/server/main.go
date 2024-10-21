package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, http.StatusOK, "server is up and running")
	})

	log.Printf("🚀 Server started on port %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, router))
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

func writeJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
