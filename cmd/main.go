package main

import (
	"log"
	"net/http"

	jsonrpcgolang "jsonrpc"
)

func main() {
	// Load configuration
	cfg, err := jsonrpcgolang.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Register the JSON-RPC handler
	http.HandleFunc("/", jsonrpcgolang.HandleJSONRPC)

	// Log server startup
	log.Printf("JSON-RPC server is running on http://%s\n", cfg.ServerAddr)

	// Start the HTTP server
	if err := http.ListenAndServe(cfg.ServerAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
