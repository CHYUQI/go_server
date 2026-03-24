package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsEvent represents an HTTP request event for metrics collection
type MetricsEvent struct {
	method   string
	path     string
	code     string
	duration float64
}

// create a channel to send metrics events
var metricChan = make(chan MetricsEvent, 100)

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

func init() {
	go func() {
		for event := range metricChan {
			httpRequestsTotal.WithLabelValues(event.method, event.path, event.code).Inc()
			httpRequestDuration.WithLabelValues(event.method, event.path, event.code).Observe(event.duration)
		}
	}()
}

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

	dur := time.Since(start).Seconds()
	if name == "err" {
		metricChan <- MetricsEvent{
			method:   "GET",
			path:     "/api/hello",
			code:     "500",
			duration: dur,
		}
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	metricChan <- MetricsEvent{
		method:   "GET",
		path:     "/api/hello",
		code:     "200",
		duration: dur,
	}
	c.JSON(200, gin.H{
		"message": "hello" + name,
	})
	// 		"error": "internal server error",
	// 	})
	// 	return
	// }

	// httpRequestsTotal.WithLabelValues("GET", "/api/hello", "200").Inc()
	// httpRequestDuration.WithLabelValues("GET", "/api/hello", "200").Observe(time.Since(start).Seconds())
	// c.JSON(200, gin.H{
	// 	"message": "hello" + name,
	// })
}
