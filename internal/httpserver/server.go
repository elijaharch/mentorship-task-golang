package httpserver

import (
	"context"
	"net/http"

	"github.com/elijaharch/mentorship-task-golang/pkg/config"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, handler http.Handler, cfg config.ServerConfig) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
