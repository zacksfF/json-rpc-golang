package jsonrpcgolang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// JSON-RPC Request Structure
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSON-RPC Response Structure
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError represents an error in the JSON-RPC response
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// SendRequest sends a JSON-RPC request to the Ethereum node provider
func SendRequest(method string, params []interface{}) (json.RawMessage, error) {
	nodeProvider := os.Getenv("NODE_PROVIDER")
	if nodeProvider == "" {
		return nil, fmt.Errorf("missing NODE_PROVIDER in .env")
	}

	requestBody := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return nil, err
	}

	log.Printf("Sending request to Ethereum node: %s", string(requestJSON))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(nodeProvider, "application/json", bytes.NewReader(requestJSON))
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return nil, err
	}

	log.Printf("Response from Ethereum node: %s", body)

	var jsonResponse JSONRPCResponse
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	if jsonResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %v", jsonResponse.Error.Message)
	}

	return jsonResponse.Result, nil
}

// Ethereum JSON-RPC Methods
func EthBlockNumber() (json.RawMessage, error) {
	return SendRequest("eth_blockNumber", nil)
}

func EthGetBlockByNumber(blockNumber string) (json.RawMessage, error) {
	return SendRequest("eth_getBlockByNumber", []interface{}{blockNumber, true})
}

func EthGetBlockByHash(blockHash string) (json.RawMessage, error) {
	return SendRequest("eth_getBlockByHash", []interface{}{blockHash, true})
}

func EthGetLogs(filter interface{}) (json.RawMessage, error) {
	return SendRequest("eth_getLogs", []interface{}{filter})
}

func EthGetBalance(address string) (json.RawMessage, error) {
	return SendRequest("eth_getBalance", []interface{}{address, "latest"})
}

func EthGetTransactionCount(address string) (json.RawMessage, error) {
	return SendRequest("eth_getTransactionCount", []interface{}{address, "latest"})
}

func EthCall(toAddress string, data string) (json.RawMessage, error) {
	params := []interface{}{map[string]interface{}{
		"to":   toAddress,
		"data": data,
	}, "latest"}
	return SendRequest("eth_call", params)
}

func EthEstimateGas(toAddress string, data string) (json.RawMessage, error) {
	params := []interface{}{map[string]interface{}{
		"to":   toAddress,
		"data": data,
	}}
	return SendRequest("eth_estimateGas", params)
}

func EthGasPrice() (json.RawMessage, error) {
	return SendRequest("eth_gasPrice", nil)
}

func EthSendRawTransaction(rawTx string) (json.RawMessage, error) {
	return SendRequest("eth_sendRawTransaction", []interface{}{rawTx})
}

func EthGetTransactionByHash(txHash string) (json.RawMessage, error) {
	return SendRequest("eth_getTransactionByHash", []interface{}{txHash})
}

func EthChainId() (json.RawMessage, error) {
	return SendRequest("eth_chainId", nil)
}
