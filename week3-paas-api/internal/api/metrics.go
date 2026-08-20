//In this file, I added two custom application metrics: one counter for HTTP requests and one histogram for request duration.
// The request counter is labeled by method, route, and status code, so I can also analyze errors and response codes.

// promauto automatically registers these collectors with Prometheus' default registry,
// while CounterVec and HistogramVec let us attach a fixed set of labels.
package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Counts API requests by method, path and HTTP status.
var httpRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "total_http_requests",
		Help: "Total number of HTTP requests.",
	},
	[]string{"method", "path", "status"},
)

// Measures how long API requests take.
var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration in seconds.",
	},
	[]string{"method", "path"},
)

// metricsMiddleware records metrics for every API request.
func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Run the real request handler.
		c.Next()

		// Do not count Prometheus scraping itself.
		if c.Request.URL.Path == "/metrics" {
			return
		}

		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		status := strconv.Itoa(c.Writer.Status())

		// Count the request.
		httpRequests.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		// Record request duration.
		httpDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(time.Since(start).Seconds())
	}
}
