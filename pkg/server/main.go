package server

import (
	"akshidas/e-com/pkg/api"
	"akshidas/e-com/pkg/db"
	"context"
	"log"
	"net/http"
)

type APIServer struct {
	Status string
	Port   string
	Store  *db.Storage
}

// Create a new server and registers routes to it

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Write([]byte("server is up and running"))
	})

	router.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Write([]byte("server is up and running"))
	})

	s.registerRoutes(router)

	wrappedRouter := NewLogger(router)
	log.Printf("🚀 Server started on port %s", s.Port)

	log.Fatal(http.ListenAndServe(s.Port, wrappedRouter))
}

func (s *APIServer) registerRoutes(r *http.ServeMux) {
	ctx := context.Background()

	// Api
	userApi := api.NewUserApi(s.Store)
	productApi := api.NewProductApi(s.Store)
	cartApi := api.NewCartApi(s.Store)
	productCategoryApi := api.NewProductCategoriesApi(s.Store)
	// Middle wares
	middlware := api.NewMiddleWare(userApi.UserService)

	// Public Routes
	r.HandleFunc("POST /users", api.RouteHandler(userApi.Create))
	r.HandleFunc("POST /login", api.RouteHandler(userApi.Login))

	// Authenticated Routes
	r.HandleFunc("GET /profile", api.RouteHandler(middlware.IsAuthenticated(ctx, userApi.GetProfile)))
	r.HandleFunc("PUT /profile", api.RouteHandler(middlware.IsAuthenticated(ctx, userApi.UpdateProfile)))
	r.HandleFunc("POST /carts", api.RouteHandler(middlware.IsAuthenticated(ctx, cartApi.Create)))
	r.HandleFunc("GET /carts", api.RouteHandler(middlware.IsAuthenticated(ctx, cartApi.GetAll)))
	r.HandleFunc("PUT /carts/{id}", api.RouteHandler(middlware.IsAuthenticated(ctx, cartApi.Update)))
	r.HandleFunc("DELETE /carts/{id}", api.RouteHandler(middlware.IsAuthenticated(ctx, cartApi.Delete)))

	// Admin Routes
	r.HandleFunc("GET /users", api.RouteHandler(middlware.IsAdmin(ctx, userApi.GetAll)))
	r.HandleFunc("GET /users/{id}", api.RouteHandler(middlware.IsAdmin(ctx, userApi.GetOne)))

	r.HandleFunc("POST /products", api.RouteHandler(middlware.IsAdmin(ctx, productApi.Create)))

	r.HandleFunc("POST /products/categories", api.RouteHandler(middlware.IsAdmin(ctx, productCategoryApi.Create)))

	r.HandleFunc("GET /products/categories", api.RouteHandler(middlware.IsAdmin(ctx, productCategoryApi.GetAll)))

	r.HandleFunc("GET /products/categories/{id}", api.RouteHandler(middlware.IsAdmin(ctx, productCategoryApi.GetOne)))

	r.HandleFunc("GET /products", api.RouteHandler(middlware.IsAdmin(ctx, productApi.GetAll)))
	r.HandleFunc("OPTIONS /products", api.RouteHandler(middlware.IsAdmin(ctx, productApi.GetAll)))
	r.HandleFunc("GET /products/{id}", api.RouteHandler(middlware.IsAdmin(ctx, productApi.GetOne)))
	r.HandleFunc("PUT /products/{id}", api.RouteHandler(middlware.IsAdmin(ctx, productApi.Update)))
	r.HandleFunc("DELETE /products/{id}", api.RouteHandler(middlware.IsAdmin(ctx, productApi.Delete)))

}
