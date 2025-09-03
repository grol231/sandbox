package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create temporary config file
	configContent := `
rabbitmq:
  host: localhost
  port: 5672
  user: guest
  password: guest
  queue: test

api:
  url: https://example.com/api
  service_id: test_service
  pass: test_pass
  source: test_source

server:
  port: 8080
  metrics_path: /metrics

logging:
  level: info
  format: json
`
	
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	tmpFile.Close()

	// Override command line args for testing
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-config=" + tmpFile.Name()}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Test RabbitMQ config
	if cfg.RabbitMQ.Host != "localhost" {
		t.Errorf("expected host 'localhost', got '%s'", cfg.RabbitMQ.Host)
	}
	if cfg.RabbitMQ.Port != 5672 {
		t.Errorf("expected port 5672, got %d", cfg.RabbitMQ.Port)
	}

	// Test connection string
	expected := "amqp://guest:guest@localhost:5672/"
	if got := cfg.RabbitMQ.ConnectionString(); got != expected {
		t.Errorf("expected connection string '%s', got '%s'", expected, got)
	}
}

func TestConnectionString(t *testing.T) {
	rmq := RabbitMQConfig{
		Host:     "example.com",
		Port:     5672,
		User:     "testuser",
		Password: "testpass",
	}

	expected := "amqp://testuser:testpass@example.com:5672/"
	if got := rmq.ConnectionString(); got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}