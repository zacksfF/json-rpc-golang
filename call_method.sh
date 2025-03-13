#!/bin/bash
# Ethereum JSON-RPC API Method Caller
# Usage: ./call_method.sh <method_name> [parameters...]

SERVER="http://localhost:8080"

# Check if a method name was provided
if [ -z "$1" ]; then
  echo "ERROR: Method name is required"
  echo "Usage: ./call_method.sh <method_name> [parameters...]"
  echo ""
  echo "Available methods:"
  echo "  eth_blockNumber                  - Get the latest block number"
  echo "  eth_chainId                      - Get the chain ID"
  echo "  eth_getBalance <address>         - Get the balance for an address"
  echo "  eth_gasPrice                     - Get current gas price"
  echo "  eth_getBlockByNumber <number>    - Get block by number (can use 'latest')"
  echo "  eth_getBlockByHash <hash>        - Get block by hash"
  echo "  eth_call <to_address> <data>     - Make a contract call"
  exit 1
fi

METHOD=$1
shift
PARAMS=()

# Process based on the method
case $METHOD in
  "eth_blockNumber")
    # No parameters needed
    PARAMS_JSON="[]"
    ;;
  
  "eth_chainId")
    # No parameters needed
    PARAMS_JSON="[]"
    ;;
  
  "eth_getBalance")
    if [ -z "$1" ]; then
      echo "ERROR: Address parameter is required for eth_getBalance"
      echo "Usage: ./call_method.sh eth_getBalance <address>"
      exit 1
    fi
    ADDRESS=$1
    PARAMS_JSON="[\"$ADDRESS\", \"latest\"]"
    ;;
  
  "eth_gasPrice")
    # No parameters needed
    PARAMS_JSON="[]"
    ;;
  
  "eth_getBlockByNumber")
    BLOCK_NUM=${1:-"latest"}
    PARAMS_JSON="[\"$BLOCK_NUM\", true]"
    ;;
  
  "eth_getBlockByHash")
    if [ -z "$1" ]; then
      echo "ERROR: Block hash parameter is required for eth_getBlockByHash"
      echo "Usage: ./call_method.sh eth_getBlockByHash <hash>"
      exit 1
    fi
    BLOCK_HASH=$1
    PARAMS_JSON="[\"$BLOCK_HASH\", true]"
    ;;
  
  "eth_call")
    if [ -z "$1" ] || [ -z "$2" ]; then
      echo "ERROR: Both to_address and data parameters are required for eth_call"
      echo "Usage: ./call_method.sh eth_call <to_address> <data>"
      exit 1
    fi
    TO_ADDRESS=$1
    DATA=$2
    PARAMS_JSON="[{\"to\":\"$TO_ADDRESS\",\"data\":\"$DATA\"}, \"latest\"]"
    ;;
  
  *)
    echo "ERROR: Unknown method '$METHOD'"
    echo "Available methods: eth_blockNumber, eth_chainId, eth_getBalance, eth_gasPrice, eth_getBlockByNumber, eth_getBlockByHash, eth_call"
    exit 1
    ;;
esac

# Create the JSON-RPC request
REQUEST="{\"jsonrpc\":\"2.0\",\"method\":\"$METHOD\",\"params\":$PARAMS_JSON,\"id\":1}"

echo "Calling method: $METHOD"
echo "Request: $REQUEST"
echo "--------------------------"

# Send the request
RESPONSE=$(curl -s -X POST $SERVER \
  -H "Content-Type: application/json" \
  -d "$REQUEST")

# Format and display the response
echo "Response:"
echo $RESPONSE | python3 -m json.tool 2>/dev/null || echo $RESPONSE