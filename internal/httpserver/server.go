package httpserver

import "net/http"

const (
	defaultAddr = ":8080"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func NewServer(...opts) *Server {
	return &Server{Addr: defaultAddr}
}
