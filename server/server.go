package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
)

type (
	Handler interface {
		Handler() (string, http.Handler)
	}

	Server struct {
		server *http.Server
		health *health
		logger *slog.Logger
	}
)

func New(cfg Config, hdrs []Handler, opt ...Option) *Server {
	var (
		mux     = http.NewServeMux()
		options = newOptions(opt...)
		health  = new(health) // переделать
	)

	mux.Handle("GET /_health", health)
	for _, c := range hdrs {
		mux.Handle(c.Handler())
	}

	handler := newMiddlewaresChain(options, mux)
	handler = http.TimeoutHandler(handler, cfg.getTimeoutHandler(), "")

	return &Server{
		server: &http.Server{
			Addr:        cfg.getAddress(),
			ReadTimeout: cfg.getTimeoutRead(),
			Handler:     handler,
		},
		logger: options.Logger,
		health: health,
	}
}

func (s *Server) Server() error {
	s.logger.Info("Starting HTTP server on " + s.server.Addr)

	lis, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	err = s.server.Serve(lis)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.health.Shutdown()
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	s.logger.Info("HTTP server successfuly stopped")

	return nil
}
