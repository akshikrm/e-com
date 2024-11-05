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

	s.registerRoutes(router)

	wrappedRouter := NewLogger(router)
	log.Printf("🚀 Server started on port %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, wrappedRouter))
}

func (s *APIServer) registerRoutes(r *http.ServeMux) {
	// Services
	userService := services.NewUserService(s.Store.(*db.PostgresStore).DB)
	productService := services.NewProductService(s.Store.(*db.PostgresStore).DB)

	// Api
	userApi := api.NewUserApi(userService)
	productApi := api.NewProductApi(productService)

	// Middle wares
	IsAdmin := api.IsAdmin(userService)

	// Public Routes
	r.HandleFunc("POST /users", api.RouteHandler(userApi.Create))
	r.HandleFunc("POST /login", api.RouteHandler(userApi.Login))

	// Authenticated Routes
	r.HandleFunc("GET /profile", api.RouteHandler(api.IsAuthenticated(userApi.GetProfile)))
	r.HandleFunc("PUT /profile", api.RouteHandler(api.IsAuthenticated(userApi.UpdateProfile)))

	// Admin Routes
	r.HandleFunc("GET /users", api.RouteHandler(IsAdmin(userApi.GetAll)))
	r.HandleFunc("GET /users/{id}", api.RouteHandler(IsAdmin(userApi.GetOne)))

	r.HandleFunc("POST /products", api.RouteHandler(IsAdmin(productApi.Create)))
	r.HandleFunc("GET /products", api.RouteHandler(IsAdmin(productApi.GetAll)))
}
