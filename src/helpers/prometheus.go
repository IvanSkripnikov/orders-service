package helpers

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsByMethodTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_by_method_total",
			Help: "Total number of HTTP requests by methods.",
		},
		[]string{"method"},
	)
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
	)
	ResponseHttpStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_http_status",
			Help: "Total number of HTTP statuses.",
		},
		[]string{"status"},
	)

	RequestLatencySummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_latency_histogram",
			Help:       "Histogram of latency HTTP requests.",
			Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
		},
	)

	RequestLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_latency_summary",
			Help:    "Summary of latency HTTP requests.",
			Buckets: []float64{.5, .95, .99},
		},
		[]string{"method"},
	)
)

func addHttpStatusCodeToPrometheus(httpStatusCode int) {
	ResponseHttpStatus.WithLabelValues(strconv.Itoa(httpStatusCode)).Inc()
}

func RegisterCommonMetrics() {
	prometheus.MustRegister(RequestsByMethodTotal, RequestsTotal, ResponseHttpStatus, RequestLatencyHistogram, RequestLatencySummary)
}
