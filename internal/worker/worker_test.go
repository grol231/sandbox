package worker

import (
	"encoding/json"
	"testing"

	"github.com/starline/rabbitmq-worker/internal/api"
	"github.com/starline/rabbitmq-worker/internal/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	apiClient := api.NewClient(&config.APIConfig{})
	
	worker := New(cfg, apiClient)
	
	if worker == nil {
		t.Fatal("expected worker to be created, got nil")
	}
	
	if worker.config != cfg {
		t.Error("expected config to be set")
	}
	
	if worker.apiClient != apiClient {
		t.Error("expected API client to be set")
	}
}

func TestMessageUnmarshal(t *testing.T) {
	jsonData := `{
		"messages": [
			{
				"recipient": "79218897127",
				"body": "StarLine код авторизации: 2652"
			}
		]
	}`

	var msgReq MessageRequest
	err := json.Unmarshal([]byte(jsonData), &msgReq)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(msgReq.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(msgReq.Messages))
	}

	msg := msgReq.Messages[0]
	expectedRecipient := "79218897127"
	expectedBody := "StarLine код авторизации: 2652"

	if msg.Recipient != expectedRecipient {
		t.Errorf("expected recipient '%s', got '%s'", expectedRecipient, msg.Recipient)
	}

	if msg.Body != expectedBody {
		t.Errorf("expected body '%s', got '%s'", expectedBody, msg.Body)
	}
}