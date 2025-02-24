package prometheus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	latencyLowrBuckets  = []float64{0.1, 0.5, 1}
	latencyHighrBuckets = []float64{0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 7.5, 10, 30, 60}
)

const (
	reqsName         = "http_requests_total"
	latencyHighrName = "http_request_duration_highr_seconds"
	latencyLowrName  = "http_request_duration_seconds"
)

type Middleware struct {
	reqs         *prometheus.CounterVec
	latencyLowr  *prometheus.HistogramVec
	latencyHighr *prometheus.HistogramVec
}

func NewGinPrometheusMiddleware(name string) gin.HandlerFunc {
	m := Middleware{
		reqs: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        reqsName,
				Help:        "How many HTTP requests processed, partitioned by status code, method, and path.",
				ConstLabels: prometheus.Labels{"service": name},
			},
			[]string{"handler", "method", "status"},
		),
		latencyHighr: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        latencyHighrName,
				Help:        "Latency with many buckets but no API-specific labels.",
				ConstLabels: prometheus.Labels{"service": name},
				Buckets:     latencyHighrBuckets,
			},
			[]string{},
		),
		latencyLowr: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        latencyLowrName,
				Help:        "Latency with only a few buckets, grouped by handler.",
				ConstLabels: prometheus.Labels{"service": name},
				Buckets:     latencyLowrBuckets,
			},
			[]string{"handler", "method", "status"},
		),
	}

	prometheus.MustRegister(m.reqs, m.latencyHighr, m.latencyLowr)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Get the actual route pattern instead of the raw URL path
		routePattern := c.FullPath()
		if routePattern == "" {
			routePattern = "unknown"
		}

		status := c.Writer.Status()
		statusGroup := string([]byte{byte('0' + status/100), 'x', 'x'})

		// Record metrics
		m.reqs.WithLabelValues(routePattern, c.Request.Method, statusGroup).Inc()
		m.latencyHighr.WithLabelValues().Observe(time.Since(start).Seconds())
		m.latencyLowr.WithLabelValues(routePattern, c.Request.Method, statusGroup).Observe(time.Since(start).Seconds())
	}
}
