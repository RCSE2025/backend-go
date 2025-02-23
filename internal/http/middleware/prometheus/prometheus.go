package prometheus

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"strings"
	"time"
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

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs         *prometheus.CounterVec
	latencyLowr  *prometheus.HistogramVec
	latencyHighr *prometheus.HistogramVec
}

// NewPatternMiddleware returns a new prometheus Middleware handler that groups requests by the chi routing pattern.
// EX: /users/{firstName} instead of /users/bob
func NewPatternMiddleware(name string, buckets ...float64) func(next http.Handler) http.Handler {
	var m Middleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": name},
		},
		//[]string{},
		[]string{"handler", "method", "status"},
	)
	prometheus.MustRegister(m.reqs)

	m.latencyHighr = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyHighrName,
		Help:        "Latency with many buckets but no API specific labels. \nMade for more accurate percentile calculations. ",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     latencyHighrBuckets,
	},
		[]string{},
	)
	prometheus.MustRegister(m.latencyHighr)

	m.latencyLowr = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyLowrName,
		Help:        "Latency with only few buckets by handler. \nMade to be only used if aggregation by handler is important. ",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     latencyLowrBuckets,
	},
		[]string{"handler", "method", "status"},
	)

	prometheus.MustRegister(m.latencyLowr)

	return m.patternHandler
}

func (c Middleware) patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		status := string(strconv.Itoa(ww.Status())[0]) + "xx"
		defer func() {
			c.reqs.WithLabelValues(routePattern, r.Method, status).Inc()
			c.latencyHighr.WithLabelValues().Observe(time.Since(start).Seconds())
			c.latencyLowr.WithLabelValues(routePattern, r.Method, status).Observe(time.Since(start).Seconds())
		}()
	}
	return http.HandlerFunc(fn)
}
