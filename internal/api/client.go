package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/starline/rabbitmq-worker/internal/config"
	"github.com/starline/rabbitmq-worker/internal/logging"
	"github.com/starline/rabbitmq-worker/internal/metrics"
)

// Client represents HTTP API client
type Client struct {
	config     *config.APIConfig
	httpClient *http.Client
}

// NewClient creates new API client
func NewClient(cfg *config.APIConfig) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends message to API endpoint
func (c *Client) SendMessage(clientID, message string) error {
	timer := prometheus.NewTimer(metrics.APIRequestDuration)
	defer timer.ObserveDuration()

	metrics.APIRequestsSent.Inc()

	// Prepare URL parameters
	params := url.Values{}
	params.Set("clientId", clientID)
	params.Set("message", message)
	params.Set("serviceId", c.config.ServiceID)
	params.Set("pass", c.config.Pass)
	params.Set("source", c.config.Source)

	// Create full URL
	fullURL := fmt.Sprintf("%s?%s", c.config.URL, params.Encode())

	logging.Debug("sending API request", logrus.Fields{
		"url":       c.config.URL,
		"client_id": clientID,
		"message":   message,
	})

	// Create POST request
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(""))
	if err != nil {
		logging.Error("failed to create API request", err, logrus.Fields{
			"client_id": clientID,
		})
		metrics.APIRequestsFailed.Inc()
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "StarLine-RabbitMQ-Worker/1.0")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logging.Error("failed to send API request", err, logrus.Fields{
			"client_id": clientID,
			"url":       c.config.URL,
		})
		metrics.APIRequestsFailed.Inc()
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("failed to read API response", err, logrus.Fields{
			"client_id":   clientID,
			"status_code": resp.StatusCode,
		})
		metrics.APIRequestsFailed.Inc()
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		logging.Error("API request failed with error status", nil, logrus.Fields{
			"client_id":    clientID,
			"status_code":  resp.StatusCode,
			"response_body": string(body),
		})
		metrics.APIRequestsFailed.Inc()
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	metrics.APIRequestsSuccess.Inc()
	logging.Info("API request sent successfully", logrus.Fields{
		"client_id":   clientID,
		"status_code": resp.StatusCode,
	})

	return nil
}