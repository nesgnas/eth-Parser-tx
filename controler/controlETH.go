package controler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"regexp"
	"serverETH/dataStruct"
	"serverETH/transaction"
)

func GetOutTransaction() gin.HandlerFunc {
	return func(context *gin.Context) {

		address := common.HexToAddress(context.Param("address"))

		var temporaryAddress []dataStruct.Address
		temporaryAddress, _ = GetSubscribers()

		var found bool
		for _, addr := range temporaryAddress {
			if common.HexToAddress(addr.AddressSub) == address {
				found = true
				break
			}
		}

		if !found {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Address not found in subscribers list",
			})
			return
		}

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
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON format",
			})
		}

		// Validate Ethereum address format
		if !isValidEthereumAddress(input.AddressSub) {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Ethereum address format",
			})
			return
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

func isValidEthereumAddress(address string) bool {
	// Ethereum address regex pattern
	ethereumAddressPattern := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	return ethereumAddressPattern.MatchString(address)
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

func CheckServerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscribers []dataStruct.Address

		subscribers, _ = GetSubscribers()

		if len(subscribers) == 0 {
			c.JSON(200, gin.H{
				"message": "No addresses to monitor",
			})
		} else {
			c.JSON(200, gin.H{
				"message":    "Welcome to the server",
				"total_subs": len(subscribers),
				"monitoring": subscribers,
			})
		}
	}
}

type UnsubscribeInput struct {
	Address string `json:"address"`
}

func UnsubscribeSubscriber() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UnsubscribeInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON format",
			})
			return
		}

		// Validate Ethereum address format
		if !isValidEthereumAddress(input.Address) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Ethereum address format",
			})
			return
		}

		address := common.HexToAddress(input.Address)

		// Check if the address exists and is subscribed
		var temporaryAddress []dataStruct.Address
		temporaryAddress, _ = GetSubscribers()

		var subscriber dataStruct.Address

		var found bool
		for _, addr := range temporaryAddress {
			if common.HexToAddress(addr.AddressSub) == address {
				found = true
				subscriber = addr
				break
			}
		}

		if !found {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address not found in subscribers list",
			})
			return
		}

		// Unsubscribe the address by setting its status to false
		subscriber.Status = false
		if err := transaction.UnsubscribeAddress(subscriber); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unsubscribe address",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Address unsubscribed successfully",
		})
	}
}
