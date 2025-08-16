package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricSubsystem = "http"

	labelServerName = "http_server_name"
	labelMethod     = "http_method"
	labelPath       = "http_path"
	labelCode       = "http_code"
)

var (
	panicCounter  = mustPanicCounter()
	callHistogram = mustCallHistogram()
	callCounter   = mustCallCounter()
)

func mustCallHistogram() *prometheus.HistogramVec {
	return promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: metricSubsystem,
			Name:      "http_handling_seconds",
			Help:      "Histogram of resonse latency (seconds) of HTTP that had been application-level handled by the server.",
			Buckets:   prometheus.DefBuckets,
		}, []string{labelServerName, labelMethod, labelCode, labelPath},
	)
}

func mustCallCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: metricSubsystem,
			Name:      "http_requests_total",
			Help:      "Tracks the number of HTTP requests.",
		}, []string{labelServerName, labelMethod, labelCode, labelPath},
	)
}

func mustPanicCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: metricSubsystem,
		Name:      "http_panics_recovered_total",
		Help:      "Total number of HTTP requets recovered from internal panic.",
	}, []string{labelServerName, labelMethod, labelCode, labelPath})
}
