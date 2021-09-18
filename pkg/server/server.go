package server

import (
	"context"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(addr string, router *http.ServeMux) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return &Server{
		server: srv,
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
