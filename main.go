package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"serverETH/controler"
	"serverETH/dataStruct"
	"serverETH/router"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"serverETH/eth"
	"serverETH/transaction"
)

func main() {

	err := godotenv.Load(".env")

	// establish connection to Node ETH
	client := eth.OpenNode(os.Getenv("ETH_NODE"))

	// check available address in database
	var subscribers []dataStruct.Address
	subscribers, _ = controler.GetSubscribers()

	var addresses []common.Address
	for _, addr := range subscribers {
		addresses = append(addresses, common.HexToAddress(addr.AddressSub))
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	// Create a WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the GatherTransactionRealTime function in a separate goroutine
	go func() {
		defer wg.Done()
		// server always run this function
		transaction.GatherTransactionRealTime(client, addresses, chainID)
	}()

	// Set up the Gin server
	r := gin.New()
	r.Use(gin.Logger())

	router.Page(r)

	// Run the Gin server in a separate goroutine
	go func() {
		if err := r.Run(os.Getenv("SERVER_PORT")); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Wait indefinitely
	wg.Wait()
}
