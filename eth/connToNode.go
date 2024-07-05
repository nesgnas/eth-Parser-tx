package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func OpenNode(nodeAddress string) *ethclient.Client {
	client, err := ethclient.Dial(nodeAddress)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
