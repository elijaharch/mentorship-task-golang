package handler

import "net/http"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /calculations", h.Create)
	router.HandleFunc("PUT /calculations/{id}", h.Update)
}
