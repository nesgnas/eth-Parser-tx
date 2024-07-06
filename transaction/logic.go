package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"os"
	"serverETH/dataStruct"
	"serverETH/database"
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

var LastBlockNumber uint64 = 20241092

var CollectionOfSubscriber *mongo.Collection = database.OpenCollection(database.Client, os.Getenv("collection_name_subscriber"))

var CollectionOfTransaction *mongo.Collection = database.OpenCollection(database.Client, os.Getenv("collection_name_transaction"))

func GatherTransactionRealTime(client *ethclient.Client, addresses []common.Address, chainID *big.Int) {

	var transactions []dataStruct.TransactionDetails

	for {
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		if LastBlockNumber == 0 || header.Number.Uint64() > LastBlockNumber {
			//if LastBlockNumber == 0 || LastBlockNumber+1 > LastBlockNumber {
			LastBlockNumber = header.Number.Uint64()

			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(LastBlockNumber)))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("block Number = ", LastBlockNumber)

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

			if len(transactions) > 0 {
				storeTransaction(transactions)
			}

			//clear transaction array
			transactions = transactions[:0]

		}

		time.Sleep(15 * time.Second)

	}
	fmt.Println("outside")
}

func storeTransaction(transaction []dataStruct.TransactionDetails) {
	var documents []interface{}

	for _, tx := range transaction {

		doc := bson.M{
			"block_number": tx.BlockNumber,
			"hash":         tx.TxHash,
			"from":         tx.From,
			"to":           tx.To,
			"value":        tx.Value,
		}
		documents = append(documents, doc)
	}

	_, err := CollectionOfTransaction.InsertMany(context.Background(), documents)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Stored", len(documents), "transactions in MongoDB")
}

func QueryTransaction(address common.Address) ([]dataStruct.TransactionDetails, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from": address.Hex()},
			{"to": address.Hex()},
		},
	}

	findOptions := options.Find()

	cursor, err := CollectionOfTransaction.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var transactions []dataStruct.TransactionDetails

	for cursor.Next(context.Background()) {
		var tx dataStruct.TransactionDetails
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func UpdateSubscriptionStatus(address []dataStruct.Address) error {

	for _, address := range address {
		filter := bson.M{"address": address.AddressSub}
		update := bson.M{"$set": bson.M{"status": "true"}}

		_, err := CollectionOfSubscriber.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertAddress(address dataStruct.Address) error {

	filter := bson.M{"address": address.AddressSub}
	update := bson.M{"$set": bson.M{
		"address": address.AddressSub,
		"status":  address.Status,
	}}

	option := options.Update().SetUpsert(true)

	_, err := CollectionOfSubscriber.UpdateOne(context.Background(), filter, update, option)
	return err
}
