package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/elijaharch/mentorship-task-golang/internal/server"
	"github.com/elijaharch/mentorship-task-golang/pkg/config"
	"github.com/elijaharch/mentorship-task-golang/pkg/logger"
)

func main() {
	cfg := config.Load()

	log := logger.New(os.Stdout, cfg.Log.Level, cfg.Log.Format)
	slog.SetDefault(log)

	router := http.NewServeMux()

	srv := server.New(cfg.Server, router)

	if err := srv.Run(); err != nil {
		slog.Error("server failed", "err", err)
	}
}
