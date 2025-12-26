package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "task_manager",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "task_manager",
			Subsystem: "http",
			Name:      "request_latency_seconds",
			Help:      "HTTP request latency in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Task metrics
	tasksCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "task_manager",
			Subsystem: "tasks",
			Name:      "count",
			Help:      "Current number of tasks by status",
		},
		[]string{"status"},
	)
)

// RecordRequest records an HTTP request
func RecordRequest(method, endpoint string, status int, duration float64) {
	requestsTotal.WithLabelValues(method, endpoint, statusCode(status)).Inc()
	requestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// SetTasksCount sets the tasks count gauge by status
func SetTasksCount(status string, count float64) {
	tasksCount.WithLabelValues(status).Set(count)
}

// IncrementTasksCount increments the tasks count for a given status
func IncrementTasksCount(status string) {
	tasksCount.WithLabelValues(status).Inc()
}

// DecrementTasksCount decrements the tasks count for a given status
func DecrementTasksCount(status string) {
	tasksCount.WithLabelValues(status).Dec()
}

func statusCode(code int) string {
	switch {
	case code < 300:
		return "2xx"
	case code < 400:
		return "3xx"
	case code < 500:
		return "4xx"
	default:
		return "5xx"
	}
}
