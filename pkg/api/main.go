package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/user"
	"net/http"
)

// Registers user routes to the passed in router
func RegisterUserApi(r *http.ServeMux, db *db.PostgresStore) {
	userApi := &UserApi{
		UserService: user.NewUserService(db),
	}

	r.HandleFunc("GET /users", routeHandler(userApi.GetAll))
	r.HandleFunc("POST /users", routeHandler(userApi.Create))
	r.HandleFunc("GET /users/{id}", routeHandler(userApi.GetOne))
	r.HandleFunc("PUT /users/{id}", routeHandler(userApi.Update))
	r.HandleFunc("DELETE /users/{id}", routeHandler(userApi.Delete))
}
