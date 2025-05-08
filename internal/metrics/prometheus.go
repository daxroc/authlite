package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	AuthRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authlite_requests_total",
			Help: "Total authlite authentication requests",
		},
		[]string{"status"},
	)
	AuthDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "authlite_request_duration_seconds",
			Help:    "Authlite request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)
)

func Register() {
	prometheus.MustRegister(AuthRequests, AuthDuration)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
