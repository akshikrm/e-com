package server

import (
	"akshidas/e-com/pkg/api"
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/services"
	"log"
	"net/http"
)

type Database interface {
	Connect() error
	Init()
}

type APIServer struct {
	Status string
	Port   string
	Store  Database
}

// Create a new server and registers routes to it
func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is up and running"))
	})

	RegisterUserApi(router, s.Store)
	// RegisterProductApi(router, s.Store)

	wrappedRouter := NewLogger(router)
	log.Printf("🚀 Server started on port %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, wrappedRouter))
}

func RegisterUserApi(r *http.ServeMux, store Database) {
	userService := services.NewUserService(store.(*db.PostgresStore).DB)
	userApi := api.NewUserApi(userService)

	IsAdmin := api.IsAdmin(userService)
	r.HandleFunc("POST /users", api.RouteHandler(userApi.Create))
	r.HandleFunc("POST /login", api.RouteHandler(userApi.Login))

	r.HandleFunc("GET /profile", api.RouteHandler(api.IsAuthenticated(userApi.GetProfile)))
	r.HandleFunc("PUT /profile", api.RouteHandler(api.IsAuthenticated(userApi.UpdateProfile)))

	r.HandleFunc("GET /users", api.RouteHandler(IsAdmin(userApi.GetAll)))
	r.HandleFunc("GET /users/{id}", api.RouteHandler(IsAdmin(userApi.GetOne)))

	productService := services.NewProductService(store.(*db.PostgresStore).DB)
	productApi := api.NewProductApi(productService)
	r.HandleFunc("POST /products", api.RouteHandler(IsAdmin(productApi.Create)))
	r.HandleFunc("GET /products", api.RouteHandler(IsAdmin(productApi.GetAll)))

}

// func RegisterProductApi(r *http.ServeMux, store Database) {}
