package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// define a Prometheus counter metric
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_http_request_duration_seconds",
			Help: "Distribution of HTTP request durations in seconds",
		},
		[]string{"method", "path", "code"},
	)
)

func main() {
	r := gin.Default()

	//monitor interface
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//health check interface
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	//example interface
	r.GET("/api/hello", hellohandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

func hellohandler(c *gin.Context) {
	start := time.Now()
	name := c.DefaultQuery("name", "world")

	if name == "err" {
		httpRequestsTotal.WithLabelValues("GET", "/api/hello", "500").Inc()
		httpRequestDuration.WithLabelValues("GET", "/api/hello", "500").Observe(time.Since(start).Seconds())
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	httpRequestsTotal.WithLabelValues("GET", "/api/hello", "200").Inc()
	httpRequestDuration.WithLabelValues("GET", "/api/hello", "200").Observe(time.Since(start).Seconds())
	c.JSON(200, gin.H{
		"message": "hello" + name,
	})
}
