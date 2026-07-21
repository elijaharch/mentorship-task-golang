package main

import (
	"log"
	"net/http"

	"github.com/elijaharch/mentorship-task-golang/internal/httpserver"
	"github.com/elijaharch/mentorship-task-golang/pkg/config"
)

func main() {
	cfg := config.Load()

	router := http.NewServeMux()

	srv := httpserver.New(":"+cfg.Port, router, cfg.Server)

	if err := srv.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
