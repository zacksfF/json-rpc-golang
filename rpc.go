package jsonrpcgolang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSON-RPC Request Structure
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	//In Go, interface{} (often called the "empty interface") represents a type that can hold any value.
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSON-RPC Response Structure
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError represents an error in the JSON-RPC response
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendRequest sends a JSON-RPC request to the Ethereum node provider
func SendRequest(method string, params []interface{}) (*JSONRPCResponse, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1, // This can be dynamic or fixed
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return nil, err
	}

	resp, err := http.Post(NodeProvider, "application/json", bytes.NewReader(requestJSON))
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var jsonResponse JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	if jsonResponse.Error != nil {
		return &jsonResponse, fmt.Errorf("RPC error: %v", jsonResponse.Error.Message)
	}

	return &jsonResponse, nil
}

// eth_blockNumber - Get the latest block number
func EthBlockNumber(w http.ResponseWriter, r *http.Request) {
	response, err := SendRequest("eth_blockNumber", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling eth_blockNumber: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(response)
}

// eth_getBlockByNumber - Get a block by its number
func EthGetBlockByNumber(blockNumber string) (json.RawMessage, error) {
	params := []interface{}{blockNumber, true}
	response, err := SendRequest("eth_getBlockByNumber", params)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// eth_getBlockByHash - Get a block by its hash
func EthGetBlockByHash(blockHash string) (json.RawMessage, error) {
	params := []interface{}{blockHash, true}
	response, err := SendRequest("eth_getBlockByHash", params)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// eth_getLogs - Get logs based on filter criteria
func EthGetLogs(filter interface{}) (json.RawMessage, error) {
	params := []interface{}{filter}
	response, err := SendRequest("eth_getLogs", params)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// eth_getBalance - Get the balance of an Ethereum address
func EthGetBalance(address string) (string, error) {
	params := []interface{}{address, "latest"}
	response, err := SendRequest("eth_getBalance", params)
	if err != nil {
		return "", err
	}
	return string(response.Result), nil
}

// eth_getTransactionCount - Get the number of transactions sent from an address
func EthGetTransactionCount(address string) (string, error) {
	params := []interface{}{address, "latest"}
	response, err := SendRequest("eth_getTransactionCount", params)
	if err != nil {
		return "", err
	}
	return string(response.Result), nil
}

// eth_call - Call a contract method (not broadcasting a transaction)
func EthCall(toAddress string, data string) (json.RawMessage, error) {
	params := []interface{}{map[string]interface{}{
		"to":   toAddress,
		"data": data,
	}, "latest"}
	response, err := SendRequest("eth_call", params)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// eth_estimateGas - Estimate gas usage for a transaction
func EthEstimateGas(toAddress string, data string) (string, error) {
	params := []interface{}{map[string]interface{}{
		"to":   toAddress,
		"data": data,
	}}
	response, err := SendRequest("eth_estimateGas", params)
	if err != nil {
		return "", err
	}
	return string(response.Result), nil
}

// eth_gasPrice - Get the current gas price
func EthGasPrice() (string, error) {
	response, err := SendRequest("eth_gasPrice", nil)
	if err != nil {
		return "", err
	}
	return string(response.Result), nil
}

// eth_sendRawTransaction - Send a raw transaction
func EthSendRawTransaction(rawTx string) (string, error) {
	params := []interface{}{rawTx}
	response, err := SendRequest("eth_sendRawTransaction", params)
	if err != nil {
		return "", err
	}
	return string(response.Result), nil
}

// eth_getTransactionByHash - Get details of a transaction by its hash
func EthGetTransactionByHash(txHash string) (json.RawMessage, error) {
	params := []interface{}{txHash}
	response, err := SendRequest("eth_getTransactionByHash", params)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// eth_chainId - Get the chain ID of the current Ethereum network
func EthChainId() (string, error) {
	response, err := SendRequest("eth_chainId", nil)
	if err != nil {
		return "", err
	}
	return string(response.Result), nil
}
