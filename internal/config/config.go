package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	API      APIConfig      `yaml:"api"`
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// RabbitMQConfig holds RabbitMQ connection settings
type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Queue    string `yaml:"queue"`
}

// APIConfig holds API settings
type APIConfig struct {
	URL       string `yaml:"url"`
	ServiceID string `yaml:"service_id"`
	Pass      string `yaml:"pass"`
	Source    string `yaml:"source"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port        int    `yaml:"port"`
	MetricsPath string `yaml:"metrics_path"`
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// ConnectionString returns formatted RabbitMQ connection string
func (r *RabbitMQConfig) ConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
}

// Load reads configuration from specified file path
func Load() (*Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	flag.Parse()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}