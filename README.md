### **Project: JSON-RPC Server in Golang for Ethereum**  

This project builds a **JSON-RPC server in Golang** to interact with Ethereum blockchain nodes. It allows users to:  
✅ Fetch blockchain data (latest block, wallet balance)  
✅ Call smart contract functions (like `balanceOf` in ERC-20 tokens)  
✅ Act as a middleware between clients and Ethereum nodes  

### **Example Requests & Responses**  

1️⃣ **Get Latest Block Number**  
- **Request:** `eth_blockNumber`  
- **Response:** `"result": "0x10d4f"` (Hex block number)  

2️⃣ **Get Wallet Balance**  
- **Request:** `eth_getBalance("0x742d…", "latest")`  
- **Response:** `"result": "0x8ac7230489e80000"` (10 ETH in Wei)  

3️⃣ **Call Smart Contract (`balanceOf`)**  
- **Request:** `eth_call` with encoded `balanceOf(address)`  
- **Response:** `"result": "0x0000000000000000000000000000000000000000000000000000000000989680"` (10M USDT)  

