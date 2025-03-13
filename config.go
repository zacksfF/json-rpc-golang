package jsonrpcgolang

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration variables
type Config struct {
	ServerAddr     string
	NodeProvider   string
	RequestTimeout time.Duration
}

// LoadConfig loads the environment variables
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:8080" // Default value
	}

	nodeProvider := os.Getenv("NODE_PROVIDER")
	if nodeProvider == "" {
		return nil, err
	}

	timeoutStr := os.Getenv("REQUEST_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "10s"
	}

	requestTimeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Printf("Invalid REQUEST_TIMEOUT value, using default 10s: %v", err)
		requestTimeout = 10 * time.Second
	}

	return &Config{
		ServerAddr:     serverAddr,
		NodeProvider:   nodeProvider,
		RequestTimeout: requestTimeout,
	}, nil
}
