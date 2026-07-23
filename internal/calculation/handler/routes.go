package handler

import "net/http"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	http.HandleFunc("POST /calculations", h.Create)
}
