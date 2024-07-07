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

- **Endpoint:** `GET /blockNumber`
- **Response:**
  ```json
  {
    "currentBlock": 20241092
  }

### Subscribe/ Unsubscribe to an Address

- **Endpoint:** `POST /subscriptions` | `PUT /subscriptions`
- **Request for Subscribe (POST) :**
  ```json
  {
    "address": "0x06c2B1eB84B1147f497f11D8476A4491F73e578D"
  }
- **Responsefor Subscribe :**
  ```json
  {
    "message": "Address added successfully"
  }
- **Request for Unsubscribe (PUT) :**
  ```json
  {
    "address": "0x06c2B1eB84B1147f497f11D8476A4491F73e578D"
  }
- **Response for Unsubscribe :**
  ```json
  {
     "message": "Address unsubscribed successfully"
  }

### Get Transactions for a Subscribed Address

- **Endpoint:** `GET /transactions`
- **Query parameter:** `address`
- mercuryo : 0x8C8D7C46219D9205f056f28fee5950aD564d7465
- **Response:** Returns transactions for the specified address
  ```json
  {
  "value": [
        {
            "BlockNumber": 20246250,
            "TxHash": "0x403311e0e6b40e5e8d9da9eb71c0ce385bea09994aba0fd17a489fff612ad20d",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0x73e5361E37DFaf3fD47725474C91b00C74ac3736",
            "Value": "0.07123431689"
        },
        {
            "BlockNumber": 20246251,
            "TxHash": "0x812b34c13cbafc339957eb69fb117e760a35235d17527ff1a8f28bce54f189a5",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
            "Value": "0"
        },
        {
            "BlockNumber": 20246260,
            "TxHash": "0xa30b4c074d6b34de2f706b6defec5ead0a1e57665920551cd852cdb8f24fd65c",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0xC1c29F8773Ee94c93B3BF4321293193fEA9eA337",
            "Value": "0.03572900988"
        },
        {
            "BlockNumber": 20246260,
            "TxHash": "0x733494ff8ac75f9460e9767440f2d16a375d53d3cfb7cbe987a2ab10492a674c",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0xfb51FC3f73640E47221652FB7064B15286637e92",
            "Value": "0.0171069171"
        },
        {
            "BlockNumber": 20246265,
            "TxHash": "0x2c249c0c827317e9739f24c43a9af0ae459c59dd0ebbc90c8d0d2af6e1c61e48",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0x650b365649c42fD6e9Ca678dC5B36594E1D64542",
            "Value": "0.02684166409"
        },
        {
            "BlockNumber": 20246265,
            "TxHash": "0xc63d9103cc4f079b4dc3f242eba9c499cb173257ad083b9d105c9849fddb0474",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0x444F1F036E4002B8496124B784C1262212Fb5CcD",
            "Value": "0.1789000218"
        },
        {
            "BlockNumber": 20246265,
            "TxHash": "0x0ddb661d29dd0232c203431f4213b2ae4165e023158f3a2a4c0cf56ec30e5179",
            "From": "0x8C8D7C46219D9205f056f28fee5950aD564d7465",
            "To": "0x535c76D83A40603fdbC319b8e739E7d564FFF9c1",
            "Value": "0.01414604514"
        }
    ]
}

## SETUP

1. **Clone the repository:**
  ```bash
   git clone https://github.com/nesgnas/serverETH
   cd serverETH
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

![image](https://github.com/nesgnas/serverETH/assets/90855639/4f2c9370-0242-46e7-bb12-3a6476d87deb)


## TESTING WITH POSTMAN

1. **`GET /blockNumber`:**

   ![image](https://github.com/nesgnas/serverETH/assets/90855639/81d6e9d0-d18e-4e7c-9e69-42bb377eb539)

2. **`POST /subscriptions`:**
- Nomal case:

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/d12a6dd1-28d5-4e76-a2f5-777ee2c8217d)
- Invalid Address format

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/618a38fa-4e35-4793-9332-baa8c07de109)
- Invalid Json format

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/9e9ed714-fb82-48ef-b5f8-ccc1618558d5)

3. **`PUT /subscriptions`:**
- Normal case:

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/07c4c940-98da-450c-95e3-ee0da15fe782)
- Address not found:

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/0ecdfa8d-781c-48b3-8fca-1a79b16b9e37)
- Invalid Address format

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/47ac9189-3d74-47c4-948d-0d93d1b073e7)
- Invalid JSON format

  ![image](https://github.com/nesgnas/serverETH/assets/90855639/eb1e28a5-2359-4f4d-aba1-47907b97b286)

3. **`GET /transactions`:**

   ![image](https://github.com/nesgnas/serverETH/assets/90855639/ad1f72a7-f024-40eb-a29a-a30fc24df63a)


