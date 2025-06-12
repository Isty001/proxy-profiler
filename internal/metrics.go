package internal

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metricsItem struct {
}

type MetricsCollector struct {
	resquestDurations *prometheus.HistogramVec
}

func NewMetricsCollector(config Config) MetricsCollector {
	resquestDurations :=  prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_duration_ms",
		Help: "Request durations",
		Buckets: config.Metrics.RequestDurationMilliseconds.Buckets,
	}, []string{"method", "path", "status"})

	prometheus.MustRegister(resquestDurations)

	return MetricsCollector{
		resquestDurations: resquestDurations,
	}
}

func (s MetricsCollector) SaveRequestDuration(request *http.Request, response *http.Response, duration time.Duration) {
	s.resquestDurations.WithLabelValues(request.Method, request.URL.Path, response.Status).Observe(float64(duration / time.Millisecond))
}
