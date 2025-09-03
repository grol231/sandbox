package worker

// MessageRequest represents incoming message from RabbitMQ
type MessageRequest struct {
	Messages []Message `json:"messages"`
}

// Message represents a single message to be sent
type Message struct {
	Recipient string `json:"recipient"`
	Body      string `json:"body"`
}