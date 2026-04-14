package main

import (
	"log/slog"
	"os"

	"github.com/sheixpeer/disk-service/internal/config"
	"github.com/sheixpeer/disk-service/internal/lib/logger/sl"
	"github.com/sheixpeer/disk-service/internal/repository/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting disk-service", slog.String("env", cfg.Env))
	log.Debug("debug message are enabled")

	repo, err := postgres.New(cfg.DatabaseUrl)
	if err != nil {
		log.Error("failed to init repository", sl.Err(err))
		os.Exit(1)
	}

	_ = repo

	// TODO: init router: chi, "chi render"

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
