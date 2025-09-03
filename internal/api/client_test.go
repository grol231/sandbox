package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/starline/rabbitmq-worker/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.APIConfig{
		URL:       "https://example.com/api",
		ServiceID: "test_service",
		Pass:      "test_pass",
		Source:    "test_source",
	}

	client := NewClient(cfg)
	
	if client == nil {
		t.Fatal("expected client to be created, got nil")
	}
	
	if client.config != cfg {
		t.Error("expected config to be set")
	}
}

func TestSendMessage(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Check query parameters
		query := r.URL.Query()
		expectedParams := map[string]string{
			"clientId":  "79218897127",
			"message":   "Test message",
			"serviceId": "test_service",
			"pass":      "test_pass",
			"source":    "test_source",
		}

		for key, expectedValue := range expectedParams {
			if got := query.Get(key); got != expectedValue {
				t.Errorf("expected %s='%s', got '%s'", key, expectedValue, got)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	cfg := &config.APIConfig{
		URL:       server.URL,
		ServiceID: "test_service",
		Pass:      "test_pass",
		Source:    "test_source",
	}

	client := NewClient(cfg)
	
	err := client.SendMessage("79218897127", "Test message")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSendMessageHTTPError(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := &config.APIConfig{
		URL:       server.URL,
		ServiceID: "test_service",
		Pass:      "test_pass",
		Source:    "test_source",
	}

	client := NewClient(cfg)
	
	err := client.SendMessage("79218897127", "Test message")
	if err == nil {
		t.Error("expected error for HTTP 500 response")
	}
}