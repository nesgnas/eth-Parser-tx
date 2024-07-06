package controler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"net/http"
	"serverETH/dataStruct"
	"serverETH/eth"
	"serverETH/transaction"
)

var client *ethclient.Client
var transactions []dataStruct.TransactionDetails

func Gathering() gin.HandlerFunc {
	return func(context *gin.Context) {

		client = eth.OpenNode("https://eth-pokt.nodies.app")

		blockNumber := big.NewInt(20238323)
		address := common.HexToAddress("0xD2C82F2e5FA236E114A81173e375a73664610998")
		transactions = transaction.GatherTransaction(blockNumber, address, client)
		context.JSON(http.StatusOK, gin.H{
			"data": transactions,
		})

	}
}

func GetOutTransaction() gin.HandlerFunc {
	return func(context *gin.Context) {

		address := common.HexToAddress(context.Param("address"))

		transactions, err := transaction.QueryTransaction(address)
		if err != nil {
			log.Println("Error querying transactions:", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
			return
		}

		transactionJSON, err := json.Marshal(transactions)
		if err != nil {
			log.Println("Error serializing transactions:", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize transactions"})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"value": json.RawMessage(transactionJSON),
		})
	}
}

func GetCurrentBlockNumber() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"CurrentBlock": transaction.LastBlockNumber,
		})
	}
}

type AddressInput struct {
	AddressSub string `json:"address"`
}

func AddSubscriber() gin.HandlerFunc {
	return func(context *gin.Context) {

		var input AddressInput

		if err := context.ShouldBindJSON(&input); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		}

		newAddress := dataStruct.Address{
			AddressSub: input.AddressSub,
			Status:     true,
		}

		if err := transaction.InsertAddress(newAddress); err != nil {
			log.Println("Error inserting address:", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert address"})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Address added successfully"})
	}
}
func GetSubscribers() ([]dataStruct.Address, error) {
	filter := bson.M{"status": true}

	findOptions := options.Find()

	cursor, err := transaction.CollectionOfSubscriber.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var subscribers []dataStruct.Address

	for cursor.Next(context.Background()) {
		var subscriber dataStruct.Address
		if err := cursor.Decode(&subscriber); err != nil {
			return nil, err
		}
		fmt.Println("query ", subscriber.AddressSub, subscriber.Status)
		subscribers = append(subscribers, subscriber)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
}
