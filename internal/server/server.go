package server

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(port int, handler http.Handler, readHeaderTimeout time.Duration) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + strconv.Itoa(port),
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
