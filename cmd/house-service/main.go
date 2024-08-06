package main

import (
	"avito-backend-bootcamp/internal/config"
	"avito-backend-bootcamp/internal/http/server"
	sender "avito-backend-bootcamp/internal/infra/email"
	"avito-backend-bootcamp/internal/infra/jwt"
	"avito-backend-bootcamp/internal/infra/repository/postgres"
	"avito-backend-bootcamp/internal/service/auth"
	emailSender "avito-backend-bootcamp/internal/service/email-sender"
	"avito-backend-bootcamp/internal/service/flat"
	"avito-backend-bootcamp/internal/service/house"
	sub "avito-backend-bootcamp/internal/service/subscription"
	"time"

	"avito-backend-bootcamp/pkg/utils/flags"
	"avito-backend-bootcamp/pkg/utils/sl"

	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// Parse flags
	flags := flags.MustParseFlags()

	// Load config
	cfg := config.MustLoad(flags.ConfigPath)

	// Initialize logger
	log := setupLogger(cfg.Env)
	log.Info("starting house-service", slog.String("env", cfg.Env))

	// Initialize dependencies
	validate := validator.New()
	jwtManager := jwt.New(cfg.JWT.SecretKey, cfg.JWT.TokenTTL)
	emailClient := sender.New()
	repository, err := postgres.New(context.Background(), &cfg.DB)
	if err != nil {
		log.Error("failed to create DB", sl.Err(err))
		os.Exit(1)
	}

	// Initialize services
	flatService := flat.New(log, repository, repository)
	houseService := house.New(log, repository)
	subService := sub.New(log, repository)
	authService := auth.New(log, jwtManager, repository)
	emailService := emailSender.New(log, emailClient, repository, repository, repository)

	// Initialize server
	srv, err := server.New(
		cfg, log,
		validate,
		authService,
		flatService,
		houseService,
		subService,
		jwtManager,
	)
	if err != nil {
		log.Error("failed to create server", sl.Err(err))
		os.Exit(1)
	}

	// Start background worker
	emailService.StartProcessEvents(context.Background(), 5*time.Second)

	// Start server
	go func() {
		err := srv.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server")
		}
	}()

	log.Info("starting server", slog.String("address", cfg.Address))

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Info("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully stop server", sl.Err(err))
		return
	}
	log.Info("server stopped gracefully")
}

func setupLogger(env string) (log *slog.Logger) {
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return
}
