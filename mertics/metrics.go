package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	})

	ResponseTimeHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "HTTP response time in seconds.",
		Buckets: prometheus.ExponentialBuckets(0.01, 2, 10),
	})

    ThroughputHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "http_throughput_bytes_per_second",
		Help: "HTTP throughput bytes per second.",
        Buckets: prometheus.ExponentialBuckets(100,2,10),
	})

	HttpStatusCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_count",
			Help: "Number of responses per status code.",
		},
		[]string{"status_code"},
	)

	ErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of proxy errors.",
	})
)