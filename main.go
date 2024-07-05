package main

import (
	"context"
	"fmt"
	"log"
	"serverETH/router"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"serverETH/eth"
	"serverETH/transaction"
)

func main() {
	client := eth.OpenNode("https://eth-pokt.nodies.app")

	addresses := []common.Address{
		common.HexToAddress("0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97"),
		// Add more addresses as needed
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	fmt.Println("was heheheheeh")

	// Create a WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the GatherTransctionRealTime function in a separate goroutine
	go func() {
		defer wg.Done()
		transaction.GatherTransctionRealTime(client, addresses, chainID)
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
