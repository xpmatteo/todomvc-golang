package web

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		var body interface{}
		if len(r.Form) > 0 {
			body = r.Form
		} else {
			body = ""
		}
		log.Printf("%-4s %s %s", r.Method, r.URL.String(), body)
	})
}

func Slowdown(delayMilli int, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(delayMilli) * time.Millisecond)
		handler.ServeHTTP(w, r)
	})
}

// Metrics generates Prometheus metrics for a handler
//
// The implementation is borrowed from the [Prometheus examples]
//
// [Prometheus examples]: https://github.com/prometheus/client_golang/blob/v1.18.0/examples/middleware/httpmiddleware/httpmiddleware.go
func Metrics(handlerName string, handler http.Handler) http.Handler {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": handlerName}, prometheus.DefaultRegisterer)

	totalOpts := prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Tracks the number of HTTP requests.",
	}
	requestsTotal := promauto.With(reg).NewCounterVec(totalOpts, []string{"method", "code"})

	const bucketStart = 0.1
	const bucketFactor = 1.5
	const bucketsCount = 5
	durationOpts := prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Tracks the latencies for HTTP requests.",
		Buckets: prometheus.ExponentialBuckets(bucketStart, bucketFactor, bucketsCount),
	}
	requestDuration := promauto.With(reg).NewHistogramVec(durationOpts, []string{"method", "code"})

	return promhttp.InstrumentHandlerCounter(requestsTotal,
		promhttp.InstrumentHandlerDuration(requestDuration,
			handler))
}
