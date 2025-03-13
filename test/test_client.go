package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func main() {
	fmt.Println("Testing JSON-RPC Ethereum Server")

	// Test requests
	testRequests := []JSONRPCRequest{
		{
			JSONRPC: "2.0",
			Method:  "eth_blockNumber",
			Params:  []interface{}{},
			ID:      1,
		},
		{
			JSONRPC: "2.0",
			Method:  "eth_chainId",
			Params:  []interface{}{},
			ID:      2,
		},
		{
			JSONRPC: "2.0",
			Method:  "eth_getBalance",
			Params:  []interface{}{"0x742d35Cc6634C0532925a3b844Bc454e4438f44e"},
			ID:      3,
		},
		{
			JSONRPC: "2.0",
			Method:  "eth_gasPrice",
			Params:  []interface{}{},
			ID:      4,
		},
	}

	for _, req := range testRequests {
		testRequest(req)
	}
}

func testRequest(request JSONRPCRequest) {
	requestJSON, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
		return
	}

	resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewReader(requestJSON))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Format the JSON response for better readability
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}

	fmt.Printf("Method: %s\nResponse:\n%s\n\n", request.Method, prettyJSON.String())
}
