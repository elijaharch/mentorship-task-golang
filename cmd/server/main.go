package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/elijaharch/mentorship-task-golang/internal/config"
	"github.com/elijaharch/mentorship-task-golang/internal/db"
	"github.com/elijaharch/mentorship-task-golang/internal/handler"
	"github.com/elijaharch/mentorship-task-golang/internal/logger"
	"github.com/elijaharch/mentorship-task-golang/internal/repository"
	"github.com/elijaharch/mentorship-task-golang/internal/server"
	"github.com/elijaharch/mentorship-task-golang/internal/service"
	"github.com/elijaharch/mentorship-task-golang/migrations"
)

func main() {
	if err := run(); err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	log := logger.New(os.Stdout, cfg.Log.Level, cfg.Log.Format)
	slog.SetDefault(log)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	database, err := db.Open(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer database.Close()

	log.Info("applying migrations")
	if err := db.Migrate(ctx, database, migrations.FS); err != nil {
		return err
	}

	calcRepo := repository.New(database)
	calcSvc := service.New(calcRepo)

	router := server.NewRouter(server.Handler{
		Calculation: handler.New(calcSvc, log),
	}, log)

	srv := server.New(cfg.Server, router, log)

	slog.Info("starting application")
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	slog.Info("application stopped")
	return nil
}
