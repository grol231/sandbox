package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestStartMetricsServer(t *testing.T) {
	// Test metrics endpoint
	handler := promhttp.Handler()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	
	body := w.Body.String()
	if !strings.Contains(body, "# HELP") {
		t.Error("expected Prometheus metrics format")
	}
}

func TestMetricsIncrement(t *testing.T) {
	// Reset metrics for clean test
	MessagesReceived.Add(0)
	MessagesProcessed.Add(0)
	APIRequestsSent.Add(0)
	
	// Test incrementing counters
	MessagesReceived.Inc()
	MessagesProcessed.Inc()
	APIRequestsSent.Inc()
	
	// Note: In real tests, you might want to use prometheus testutil package
	// to properly test metric values. This is a basic structure test.
}

func TestWorkerHealthyGauge(t *testing.T) {
	// Test setting gauge values
	WorkerHealthy.Set(1)
	WorkerHealthy.Set(0)
	
	// In a real test, you'd verify the actual gauge value
	// using prometheus testutil package
}

func TestDurationHistogram(t *testing.T) {
	// Test histogram timer
	timer := prometheus.NewTimer(MessageProcessingDuration)
	time.Sleep(1 * time.Millisecond) // Simulate work
	timer.ObserveDuration()
	
	timer2 := prometheus.NewTimer(APIRequestDuration)
	time.Sleep(1 * time.Millisecond) // Simulate work
	timer2.ObserveDuration()
}