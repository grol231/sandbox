package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInit(t *testing.T) {
	tests := []struct {
		level    string
		format   string
		expected logrus.Level
	}{
		{"debug", "json", logrus.DebugLevel},
		{"info", "json", logrus.InfoLevel},
		{"warn", "json", logrus.WarnLevel},
		{"error", "json", logrus.ErrorLevel},
		{"invalid", "json", logrus.InfoLevel}, // default fallback
	}

	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			logger := Init(test.level, test.format)
			
			if logger.Level != test.expected {
				t.Errorf("expected level %v, got %v", test.expected, logger.Level)
			}

			if test.format == "json" {
				if _, ok := logger.Formatter.(*logrus.JSONFormatter); !ok {
					t.Error("expected JSON formatter")
				}
			} else {
				if _, ok := logger.Formatter.(*logrus.TextFormatter); !ok {
					t.Error("expected text formatter")
				}
			}
		})
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := Init("info", "json")
	logger.SetOutput(&buf)

	Info("test message", logrus.Fields{"key": "value"})

	output := buf.String()
	
	// Parse JSON to verify structure
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("failed to parse log output as JSON: %v", err)
	}

	if logEntry["msg"] != "test message" {
		t.Errorf("expected message 'test message', got '%v'", logEntry["msg"])
	}

	if logEntry["level"] != "info" {
		t.Errorf("expected level 'info', got '%v'", logEntry["level"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("expected field 'key'='value', got '%v'", logEntry["key"])
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	logger := Init("info", "json")
	logger.SetOutput(&buf)

	testErr := errors.New("test error")
	Error("test error message", testErr, logrus.Fields{"error_type": "test"})

	output := buf.String()
	
	if !strings.Contains(output, "test error message") {
		t.Error("expected log output to contain error message")
	}

	if !strings.Contains(output, "error_type") {
		t.Error("expected log output to contain custom field")
	}
}