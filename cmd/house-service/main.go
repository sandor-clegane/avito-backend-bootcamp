package main

import (
	"avito-backend-bootcamp/internal/config"
	"avito-backend-bootcamp/internal/di"
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

	// Initialize dependencies using DI
	di := di.New(cfg, log)

	// Start background worker
	di.GetSenderService().StartProcessEvents(context.Background(), 30*time.Second)

	// Start server
	go func() {
		err := di.GetHTTPServer().Run()
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

	if err := di.GetHTTPServer().Shutdown(ctx); err != nil {
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
