package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// MessagesReceived counts total messages received from RabbitMQ
	MessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_messages_received_total",
		Help: "The total number of messages received from RabbitMQ",
	})

	// MessagesProcessed counts total messages processed successfully
	MessagesProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "messages_processed_total",
		Help: "The total number of messages processed successfully",
	})

	// APIRequestsSent counts total API requests sent
	APIRequestsSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_requests_sent_total",
		Help: "The total number of API requests sent",
	})

	// APIRequestsSuccess counts successful API requests
	APIRequestsSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_requests_success_total",
		Help: "The total number of successful API requests",
	})

	// APIRequestsFailed counts failed API requests
	APIRequestsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_requests_failed_total",
		Help: "The total number of failed API requests",
	})

	// MessageProcessingDuration tracks message processing time
	MessageProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "message_processing_duration_seconds",
		Help: "Duration of message processing in seconds",
		Buckets: prometheus.DefBuckets,
	})

	// APIRequestDuration tracks API request duration
	APIRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "api_request_duration_seconds",
		Help: "Duration of API requests in seconds",
		Buckets: prometheus.DefBuckets,
	})

	// WorkerHealthy indicates if worker is healthy (1) or not (0)
	WorkerHealthy = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "worker_healthy",
		Help: "Worker health status: 1 for healthy, 0 for unhealthy",
	})
)

// StartMetricsServer starts the Prometheus metrics HTTP server
func StartMetricsServer(port string, metricsPath string) {
	http.Handle(metricsPath, promhttp.Handler())
	
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic("Failed to start metrics server: " + err.Error())
		}
	}()
}