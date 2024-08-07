package server

import (
	createFlat "avito-backend-bootcamp/internal/http/handlers/create-flat"
	createHouse "avito-backend-bootcamp/internal/http/handlers/create-house"
	dummyLogin "avito-backend-bootcamp/internal/http/handlers/dummy-login"
	getHouse "avito-backend-bootcamp/internal/http/handlers/get-house"
	login "avito-backend-bootcamp/internal/http/handlers/login"
	signup "avito-backend-bootcamp/internal/http/handlers/signup"
	subscribe "avito-backend-bootcamp/internal/http/handlers/subscribe"
	updateFlat "avito-backend-bootcamp/internal/http/handlers/update-flat"

	"avito-backend-bootcamp/internal/config"
	mwr "avito-backend-bootcamp/internal/http/middleware"
	"avito-backend-bootcamp/internal/infra/jwt"
	"avito-backend-bootcamp/internal/service/auth"
	"avito-backend-bootcamp/internal/service/flat"
	"avito-backend-bootcamp/internal/service/house"
	sub "avito-backend-bootcamp/internal/service/subscription"

	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	server *http.Server
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	validate *validator.Validate,
	authService *auth.Service,
	flatService *flat.Service,
	houseService *house.Service,
	subService *sub.Service,
	jwtManager *jwt.Manager,
) (*Server, error) {
	// init router
	router := chi.NewRouter()

	// use built-in middleware to log requests.
	router.Use(mwr.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Доступно всем, авторизация не нужна
	router.Get("/dummyLogin", dummyLogin.New(log, authService))
	router.Post("/login", login.New(log, validate, authService))
	router.Post("/signup", signup.New(log, validate, authService))

	// Доступно любому авторизированному
	router.Group(func(r chi.Router) {
		r.Use(mwr.NewAuthModeratorOrClient(jwtManager))
		r.Get("/house/{id}", getHouse.New(log, flatService))
		r.Post("/house/{id}/subscribe", subscribe.New(log, validate, subService))
		r.Post("/flat/create", createFlat.New(log, validate, flatService))
	})

	// Доступно только для модераторов
	router.Group(func(r chi.Router) {
		r.Use(mwr.NewAuthModerator(jwtManager))
		r.Post("/house/create", createHouse.New(log, validate, houseService))
		r.Post("/flat/update", updateFlat.New(log, validate, flatService))
	})

	return &Server{
		server: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		},
	}, nil
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
