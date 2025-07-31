package server

import (
	"context"
	"net/http"
	"strconv"
)

type Server struct {
	httpServer *http.Server
}

func New(port int, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
