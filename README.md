# Ethereum Blockchain Parser and HTTP API

This project implements a blockchain parser in Go for monitoring Ethereum blockchain transactions of specific addresses. It includes an HTTP API to interact with the parser and gather transaction data.

## Objective

The goal is to create a functional blockchain parser that can:

- Monitor specified Ethereum addresses for transactions.
- Collect transactions related to monitored addresses.
- Expose an HTTP API for interacting with the parsed data.

## Task

### Server Implementation

The blockchain parser does the following:

1. **Monitor Specified Addresses:** It watches specified Ethereum addresses for any new transactions.

2. **Collect Transactions:** For each monitored address, it gathers transactions starting from the subscribed date without needing historical data.

3. **HTTP API Exposure:** Provides endpoints to:
   - Get the current block number.
   - Subscribe to an Ethereum address.
   - Get transactions for a subscribed address.

### Others

- **Storage:** Uses MongoDB for storing subscribed addresses and their transactions.

- **Setup Instructions:** Includes Docker Compose setup for MongoDB.

## APIs

### Get Current Block Number

- **Endpoint:** `GET /currentBlockNumber`
- **Response:**
  ```json
  {
    "currentBlock": 20241092
  }

### Subscribe to an Address

- **Endpoint:** `POST /subscribe`
- **Request :**
  ```json
  {
    "address": "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97"
  }
- **Response:**
  ```json
  {
    "message": "Subscribed to address successfully"
  }

### Get Transactions for a Subscribed Address

- **Endpoint:** `GET /transactions`
- **Query parameter:** `address`
- **Response:** Returns transactions for the specified address
  ```json
  {
  "transactions": [
    {
      "txHash": "0x123456...",
      "from": "0xabc...",
      "to": "0xdef...",
      "value": "0.5 ETH"
    },
    {
      "txHash": "0x7890ab...",
      "from": "0xghi...",
      "to": "0xjkl...",
      "value": "1.2 ETH"
    }
  ]
}

## SETUP

1. **Clone the repository:**
  ```bash
   git clone https://github.com/your/repository.git
   cd repository
  ```

2. **Start Docker Compose for MongoDB:**
  ```bash
   docker-compose up -d
  ```
3. **Run the application:**
  ```bash
   go run main.go
  ```
The server will start at http://localhost:8000.
