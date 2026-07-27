package server

import (
	"log/slog"
	"net/http"

	"github.com/elijaharch/mentorship-task-golang/internal/features/calculation/handler"
)

type Handler struct {
	Calculation *handler.Handler
}

func NewRouter(h Handler, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /calculations", h.Calculation.Create)
	r.HandleFunc("GET /calculations/{id}", h.Calculation.Get)
	r.HandleFunc("PUT /calculations/{id}", h.Calculation.Update)

	return r
}
