package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"serverETH/dataStruct"
	"serverETH/eth"
	"time"
)

func GatherTransaction(blockNumber *big.Int, address common.Address, client *ethclient.Client) []dataStruct.TransactionDetails {

	var transactions []dataStruct.TransactionDetails

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("Failed to retrieve block: %v", err)
	}

	countTx := 0

	for _, tx := range block.Transactions() {
		countTx++

		signer := types.LatestSignerForChainID(chainID)
		from, err := types.Sender(signer, tx)
		if err == nil {
			//fmt.Println("Sender: ", from.Hex())
		} else {
			fmt.Println("error:", err)
		}

		if tx.To() != nil && (*tx.To() == address || from == address) {

			transactions = append(transactions, dataStruct.TransactionDetails{
				BlockNumber: block.Number().Uint64(),
				TxHash:      tx.Hash().Hex(),
				From:        from.Hex(),
				To:          tx.To().Hex(),
				Value:       tx.Value().String(),
			})
		}
	}
	return transactions

}

func GatherTransctionRealTime(client *ethclient.Client, addresses []common.Address, chainID *big.Int) {
	fmt.Println("was heheheheeh")
	var lastBlockNumber uint64 = 20241092

	var transactions []dataStruct.TransactionDetails
	fmt.Println("was here")

	//var counterBlock = 0

	for {
		//if counterBlock == 5 {
		//	break
		//}
		//
		//counterBlock++
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		if lastBlockNumber == 0 || header.Number.Uint64() > lastBlockNumber {
			//if lastBlockNumber == 0 || lastBlockNumber+1 > lastBlockNumber {
			lastBlockNumber = header.Number.Uint64()

			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(lastBlockNumber)))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("block Number = ", lastBlockNumber)

			count := 0
			for _, tx := range block.Transactions() {

				signed := types.LatestSignerForChainID(chainID)
				from, err := types.Sender(signed, tx)
				if err != nil {
					log.Fatal(err)
					continue
				}

				// ghet
				for _, address := range addresses {
					if tx.To() != nil && (*tx.To() == address || from == address) {
						count++
						transactions = append(transactions, dataStruct.TransactionDetails{
							BlockNumber: block.Number().Uint64(),
							TxHash:      tx.Hash().Hex(),
							From:        from.Hex(),
							To:          tx.To().Hex(),
							Value:       eth.ConvertWei(tx.Value()),
						})
					}
				}

			}
			fmt.Println("total tx caught", count)
			//store in to database
			for _, tx := range transactions {
				dataStruct.TrasactionGlobal = append(dataStruct.TrasactionGlobal, tx)
			}

			//clear transaction array
			transactions = transactions[:0]

		}

		//for _, txDetails := range transactions {
		//	fmt.Printf("Block Number: %d\n", txDetails.BlockNumber)
		//	fmt.Printf("Transaction: %s\n", txDetails.TxHash)
		//	fmt.Printf("From: %s\n", txDetails.From)
		//	fmt.Printf("To: %s\n", txDetails.To)
		//	fmt.Printf("Value: %s\n", txDetails.Value)
		//}

		time.Sleep(12 * time.Second)
		//time.Sleep(12 * time.Second)

	}
	fmt.Println("outside")

}
