package server

import (
	"log/slog"
	"net/http"

	"github.com/elijaharch/mentorship-task-golang/internal/feature/calculation/handler"
)

type Handlers struct {
	Calculation *handler.Handler
}

func NewRouter(h Handlers, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /calculations", h.Calculation.Create)
	r.HandleFunc("GET /calculations/{id}", h.Calculation.Get)
	r.HandleFunc("PUT /calculations/{id}", h.Calculation.Update)
	r.HandleFunc("DELETE /calculations/{id}", h.Calculation.Delete)

	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return chain(r, loggingMiddleware(logger))
}
