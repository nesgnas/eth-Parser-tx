package main

import (
	"context"
	"fmt"
	"log"
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

	client := eth.OpenNode("https://eth-pokt.nodies.app")

	var subscribers []dataStruct.Address

	subscribers, _ = controler.GetSubscribers()

	fmt.Println("sub ", subscribers)

	var addresses []common.Address
	for _, addr := range subscribers {
		addresses = append(addresses, common.HexToAddress(addr.AddressSub))
	}

	//addresses := []common.Address{
	//	common.HexToAddress("{\n  \"address\": \"0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97\"\n}"),
	//	// Add more addresses as needed
	//}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	fmt.Println("was heheheheeh")
	fmt.Println(addresses)

	// Create a WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the GatherTransactionRealTime function in a separate goroutine
	go func() {
		defer wg.Done()
		transaction.GatherTransactionRealTime(client, addresses, chainID)
	}()

	// Set up the Gin server
	r := gin.New()
	r.Use(gin.Logger())

	router.Page(r)

	// Define a basic endpoint
	r.GET("/landingPage", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"value": "welcome to server",
		})

	})

	// Run the Gin server in a separate goroutine
	go func() {
		if err := r.Run(":8000"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Wait indefinitely
	wg.Wait()
}
