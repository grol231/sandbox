package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	
	"github.com/starline/rabbitmq-worker/internal/api"
	"github.com/starline/rabbitmq-worker/internal/config"
	"github.com/starline/rabbitmq-worker/internal/logging"
	"github.com/starline/rabbitmq-worker/internal/metrics"
)

// Worker represents the main worker that processes RabbitMQ messages
type Worker struct {
	config    *config.Config
	apiClient *api.Client
	conn      *amqp.Connection
	channel   *amqp.Channel
}

// New creates a new worker instance
func New(cfg *config.Config, apiClient *api.Client) *Worker {
	return &Worker{
		config:    cfg,
		apiClient: apiClient,
	}
}

// Start initializes connection to RabbitMQ and starts consuming messages
func (w *Worker) Start(ctx context.Context) error {
	if err := w.connect(); err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	logging.Info("worker started successfully", logrus.Fields{
		"queue":    w.config.RabbitMQ.Queue,
		"host":     w.config.RabbitMQ.Host,
		"port":     w.config.RabbitMQ.Port,
	})

	metrics.WorkerHealthy.Set(1)

	return w.consume(ctx)
}

// connect establishes connection to RabbitMQ
func (w *Worker) connect() error {
	conn, err := amqp.Dial(w.config.RabbitMQ.ConnectionString())
	if err != nil {
		logging.Error("failed to connect to RabbitMQ", err, logrus.Fields{
			"host": w.config.RabbitMQ.Host,
			"port": w.config.RabbitMQ.Port,
		})
		metrics.WorkerHealthy.Set(0)
		return err
	}
	w.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		logging.Error("failed to open RabbitMQ channel", err)
		metrics.WorkerHealthy.Set(0)
		return err
	}
	w.channel = ch

	// Declare queue (create if not exists)
	_, err = ch.QueueDeclare(
		w.config.RabbitMQ.Queue, // queue name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		logging.Error("failed to declare queue", err, logrus.Fields{
			"queue": w.config.RabbitMQ.Queue,
		})
		metrics.WorkerHealthy.Set(0)
		return err
	}

	logging.Info("successfully connected to RabbitMQ", logrus.Fields{
		"queue": w.config.RabbitMQ.Queue,
	})

	return nil
}

// consume starts consuming messages from RabbitMQ queue
func (w *Worker) consume(ctx context.Context) error {
	msgs, err := w.channel.Consume(
		w.config.RabbitMQ.Queue, // queue
		"",                      // consumer
		false,                   // auto-ack (manual ack for reliability)
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		logging.Error("failed to register consumer", err)
		metrics.WorkerHealthy.Set(0)
		return err
	}

	logging.Info("starting message consumption", logrus.Fields{
		"queue": w.config.RabbitMQ.Queue,
	})

	for {
		select {
		case <-ctx.Done():
			logging.Info("worker context cancelled, shutting down")
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				logging.Warn("message channel closed, attempting to reconnect")
				metrics.WorkerHealthy.Set(0)
				if err := w.reconnect(); err != nil {
					return err
				}
				return w.consume(ctx) // Restart consumption
			}
			
			if err := w.processMessage(d); err != nil {
				logging.Error("failed to process message", err, logrus.Fields{
					"message_id": d.MessageId,
					"body":       string(d.Body),
				})
				// Reject message and don't requeue to prevent infinite loops
				d.Nack(false, false)
			} else {
				// Acknowledge successful processing
				d.Ack(false)
			}
		}
	}
}

// processMessage processes a single message from RabbitMQ
func (w *Worker) processMessage(delivery amqp.Delivery) error {
	timer := prometheus.NewTimer(metrics.MessageProcessingDuration)
	defer timer.ObserveDuration()

	metrics.MessagesReceived.Inc()

	logging.Debug("received message from RabbitMQ", logrus.Fields{
		"message_id":   delivery.MessageId,
		"routing_key":  delivery.RoutingKey,
		"body_length": len(delivery.Body),
	})

	// Parse JSON message
	var msgReq MessageRequest
	if err := json.Unmarshal(delivery.Body, &msgReq); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Process each message in the request
	for _, msg := range msgReq.Messages {
		if err := w.apiClient.SendMessage(msg.Recipient, msg.Body); err != nil {
			return fmt.Errorf("failed to send message via API: %w", err)
		}
		
		logging.Info("message sent successfully", logrus.Fields{
			"recipient": msg.Recipient,
			"body":      msg.Body,
		})
	}

	metrics.MessagesProcessed.Inc()
	return nil
}

// reconnect attempts to reconnect to RabbitMQ
func (w *Worker) reconnect() error {
	logging.Info("attempting to reconnect to RabbitMQ")
	
	// Close existing connections
	if w.channel != nil {
		w.channel.Close()
	}
	if w.conn != nil {
		w.conn.Close()
	}

	// Wait a bit before reconnecting
	time.Sleep(5 * time.Second)

	return w.connect()
}

// Stop gracefully shuts down the worker
func (w *Worker) Stop() error {
	logging.Info("shutting down worker")
	
	metrics.WorkerHealthy.Set(0)
	
	if w.channel != nil {
		if err := w.channel.Close(); err != nil {
			logging.Error("failed to close channel", err)
		}
	}
	
	if w.conn != nil {
		if err := w.conn.Close(); err != nil {
			logging.Error("failed to close connection", err)
		}
	}
	
	logging.Info("worker shutdown complete")
	return nil
}