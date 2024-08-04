package main

import (
	"avito-backend-bootcamp/internal/config"
	"avito-backend-bootcamp/internal/http/server"
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
	flags := flags.MustParseFlags()

	cfg := config.MustLoad(flags.ConfigPath)

	log := setupLogger(cfg.Env)
	log.Info("starting house-service", slog.String("env", cfg.Env))

	validate := validator.New()

	srv, err := server.New(cfg, log, validate)
	if err != nil {
		log.Error("failed to create server", sl.Err(err))
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := srv.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server")
		}
	}()

	log.Info("starting server", slog.String("address", cfg.Address))

	<-done
	log.Info("stopping server")

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
