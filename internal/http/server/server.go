package server

import (
	"avito-backend-bootcamp/internal/config"
	h "avito-backend-bootcamp/internal/http/handlers"
	mwr "avito-backend-bootcamp/internal/http/middleware"

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

func New(cfg *config.Config, log *slog.Logger, validate *validator.Validate) (*Server, error) {
	// init router
	router := chi.NewRouter()

	// use built-in middleware to log requests.
	router.Use(mwr.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Доступно всем, авторизация не нужна
	router.Get("/dummyLogin", h.HandleDummyLogin())
	router.Post("/login", h.HandleLogin(validate))
	router.Post("/signup", h.HandleSignup(validate))

	// Доступно любому авторизированному
	router.Group(func(r chi.Router) {
		r.Use(mwr.AuthModerator)
		r.Use(mwr.AuthUser)
		r.Get("/house/{id}", h.HandleGetHouse())
		r.Post("/house/{id}/subscribe", h.HandleSubscribeHouse(validate))
		r.Post("/flat/create", h.HandleCreateFlat(validate))
	})

	// Доступно только для модераторов
	router.Group(func(r chi.Router) {
		r.Use(mwr.AuthModerator)
		r.Post("/house/create", h.HandleCreateHouse(validate))
		r.Post("/flat/update", h.HandleUpdateFlat(validate))
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
