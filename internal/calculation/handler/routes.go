package handler

import "net/http"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /calculations", h.Create)
}
