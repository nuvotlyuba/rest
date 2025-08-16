package server

import (
	"log/slog"
)

type (
	Options struct {
		Logger      *slog.Logger
		Middlewares []Middleware
		Name        string
	}
	Option func(s *Options)
)

func newOptions(opts ...Option) Options {
	var o = Options{
		Logger: slog.Default(),
		Name:   "http",
	}
	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithLogger(l *slog.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func WithMiddleware(m ...Middleware) Option {
	return func(o *Options) {
		o.Middlewares = append(o.Middlewares, m)
	}
}
