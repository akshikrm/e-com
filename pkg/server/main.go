package server

import (
	"akshidas/e-com/pkg/api"
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/services"
	"log"
	"net/http"
	"time"
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

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(handleToWrap http.Handler) *Logger {
	return &Logger{handleToWrap}
}

// Create a new server and registers routes to it
func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is up and running"))
	})

	RegisterUserApi(router, s.Store)
	wrappedRouter := NewLogger(router)
	log.Printf("🚀 Server started on port %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, wrappedRouter))
}

func RegisterUserApi(r *http.ServeMux, store Database) {
	userModel := model.NewUserModel(store.(*db.PostgresStore).DB)
	userService := services.NewUserService(userModel)
	userApi := api.NewUserApi(userService)

	r.HandleFunc("GET /users", api.RouteHandler(userApi.GetAll))
	r.HandleFunc("POST /users", api.RouteHandler(userApi.Create))
	r.HandleFunc("POST /login", api.RouteHandler(userApi.Login))
	r.HandleFunc("GET /users/{id}", api.RouteHandler(userApi.GetOne))
	r.HandleFunc("PUT /users/{id}", api.RouteHandler(userApi.Update))
	r.HandleFunc("DELETE /users/{id}", api.RouteHandler(userApi.Delete))
}
