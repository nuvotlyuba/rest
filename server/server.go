package server

import (
	"log/slog"
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
