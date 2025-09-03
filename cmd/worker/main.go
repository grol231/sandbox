package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/starline/rabbitmq-worker/internal/api"
	"github.com/starline/rabbitmq-worker/internal/config"
	"github.com/starline/rabbitmq-worker/internal/logging"
	"github.com/starline/rabbitmq-worker/internal/metrics"
	"github.com/starline/rabbitmq-worker/internal/worker"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging
	logger := logging.Init(cfg.Logging.Level, cfg.Logging.Format)

	logging.Info("application starting", logrus.Fields{
		"version":        "1.0.0",
		"rabbitmq_host":  cfg.RabbitMQ.Host,
		"rabbitmq_port":  cfg.RabbitMQ.Port,
		"rabbitmq_queue": cfg.RabbitMQ.Queue,
		"api_url":        cfg.API.URL,
		"metrics_port":   cfg.Server.Port,
	})

	// Start metrics server
	metrics.StartMetricsServer(strconv.Itoa(cfg.Server.Port), cfg.Server.MetricsPath)
	logging.Info("metrics server started", logrus.Fields{
		"port": cfg.Server.Port,
		"path": cfg.Server.MetricsPath,
	})

	// Create API client
	apiClient := api.NewClient(&cfg.API)

	// Create worker
	w := worker.New(cfg, apiClient)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logging.Info("shutdown signal received, stopping worker")
		cancel()
	}()

	// Start worker
	logging.Info("starting RabbitMQ worker")
	if err := w.Start(ctx); err != nil {
		if err == context.Canceled {
			logging.Info("worker stopped gracefully")
		} else {
			logging.Error("worker stopped with error", err)
			logger.Exit(1)
		}
	}

	// Cleanup
	if err := w.Stop(); err != nil {
		logging.Error("failed to stop worker cleanly", err)
	}

	logging.Info("application shutdown complete")
}