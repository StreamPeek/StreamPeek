package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	KafkaBrokers string
	Port         string
}

// Load reads configuration from environment variables.
func Load() *Config {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		KafkaBrokers: brokers,
		Port:         port,
	}
}
