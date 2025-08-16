package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Middleware func(next http.Handler) http.Handler

func newMiddlewaresChain(opts Options, handler http.Handler) http.Handler {
	middlewares := make([]Middleware, 0, len(opts.Middlewares)+3)
	middlewares = append(middlewares, middlewareMetric(opts.Name))
	middlewares = append(middlewares, middlewareLoggerInit)
	middlewares = append(middlewares, middlewareRecovery(opts.Logger, opts.Name))
	middlewares = append(middlewares, opts.Middlewares...)

	for i := range middlewares {
		handler = middlewares[len(middlewares)-1-i](handler)
	}

	return handler
}

func middlewareLoggerInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")

		attrs := []slog.Attr{
			slog.String("http.method", r.Method),
			slog.String("http.path", r.URL.Path),
			slog.String("http.proto", r.Proto),
		}
		if traceID != "" {
			attrs = append(attrs, slog.String("trace_id", traceID))
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "logger_attrs", attrs)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func middlewareRecovery(l *slog.Logger, serverName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				p := recover()
				if p == nil {
					return
				}
				var (
					stack = make([]byte, 64<<10)
					msg   = fmt.Sprint("%v", p)
				)

				stack = stack[:runtime.Stack(stack, false)]
				l.ErrorContext(r.Context(), msg, slog.String("stack", string(stack)))
				w.WriteHeader(http.StatusInternalServerError)
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func middlewareMetric(serverName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				start = time.Now()
				sw    = newResponseWriter(w)
			)

			next.ServeHTTP(sw, r)

			labels := prometheus.Labels{
				labelServerName: serverName,
				labelCode:       strconv.Itoa(sw.code),
				labelMethod:     r.Method,
				labelPath:       r.URL.Path,
			}

			callCounter.With(labels).Inc()
			callHistogram.With(labels).Observe(time.Since(start).Seconds())
		})
	}
}
