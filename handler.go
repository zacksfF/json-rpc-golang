package jsonrpcgolang

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleJSONRPC dynamically routes JSON-RPC requests to the correct method
func HandleJSONRPC(w http.ResponseWriter, r *http.Request) {
	// Only handle POST requests
	if r.Method != http.MethodPost {
		http.Error(w, `{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"}}`, http.StatusMethodNotAllowed)
		return
	}

	var request JSONRPCRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	log.Printf("Received JSON-RPC request: %+v", request)

	var result json.RawMessage
	var errResp *RPCError

	switch request.Method {
	case "eth_blockNumber":
		result, errResp = processRequest(EthBlockNumber())
	case "eth_chainId":
		result, errResp = processRequest(EthChainId())
	case "eth_getBlockByNumber":
		if len(request.Params) < 1 {
			errResp = &RPCError{Code: -32602, Message: "Missing block number parameter"}
		} else {
			blockNumber, ok := request.Params[0].(string)
			if ok {
				result, errResp = processRequest(EthGetBlockByNumber(blockNumber))
			} else {
				errResp = &RPCError{Code: -32602, Message: "Invalid block number format"}
			}
		}
	case "eth_getBlockByHash":
		if len(request.Params) < 1 {
			errResp = &RPCError{Code: -32602, Message: "Missing block hash parameter"}
		} else {
			blockHash, ok := request.Params[0].(string)
			if ok {
				result, errResp = processRequest(EthGetBlockByHash(blockHash))
			} else {
				errResp = &RPCError{Code: -32602, Message: "Invalid block hash format"}
			}
		}
	case "eth_getBalance":
		if len(request.Params) < 1 {
			errResp = &RPCError{Code: -32602, Message: "Missing address parameter"}
		} else {
			address, ok := request.Params[0].(string)
			if ok {
				result, errResp = processRequest(EthGetBalance(address))
			} else {
				errResp = &RPCError{Code: -32602, Message: "Invalid address format"}
			}
		}
	case "eth_gasPrice":
		result, errResp = processRequest(EthGasPrice())
	case "eth_call":
		if len(request.Params) < 1 {
			errResp = &RPCError{Code: -32602, Message: "Missing call parameters"}
		} else {
			callData, ok := request.Params[0].(map[string]interface{})
			if ok {
				to, toOk := callData["to"].(string)
				data, dataOk := callData["data"].(string)
				if toOk && dataOk {
					result, errResp = processRequest(EthCall(to, data))
				} else {
					errResp = &RPCError{Code: -32602, Message: "Invalid call parameters"}
				}
			} else {
				errResp = &RPCError{Code: -32602, Message: "Invalid call format"}
			}
		}
	default:
		errResp = &RPCError{Code: -32601, Message: "Method not found"}
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
		Error:   errResp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// processRequest handles the response and formats errors properly
func processRequest(res json.RawMessage, err error) (json.RawMessage, *RPCError) {
	if err != nil {
		return nil, &RPCError{Code: -32000, Message: err.Error()}
	}
	return res, nil
}
