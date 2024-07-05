package controler

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
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
		fmt.Println("dataStruct.TrasactionGlobal = ", dataStruct.TrasactionGlobal)
		context.JSON(http.StatusOK, gin.H{
			"value": dataStruct.TrasactionGlobal,
		})
	}
}
