// Package server wires the HTTP server lifecycle.
package server

import (
	"context"
	"net/http"

	"go-note-api/internal/handler"
)

// Server wraps net/http.Server with graceful shutdown support.
type Server struct {
	httpSrv *http.Server
}

// New creates a Server that listens on addr and serves the given handler.
func New(addr string, h *handler.Handler) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:    addr,
			Handler: h.Routes(),
		},
	}
}

// Start begins serving requests.
func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

// Shutdown stops the server, allowing in-flight requests to finish.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
