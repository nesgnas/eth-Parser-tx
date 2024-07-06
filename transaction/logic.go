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

var LastBlockNumber uint64 = 20241092

// "CollectionOfSubscriber" declare connection of Subscriber collection
var CollectionOfSubscriber *mongo.Collection = database.OpenCollection(database.Client, os.Getenv("collection_name_subscriber"))

// "CollectionOfTransaction" declare connection of Transaction complexion
var CollectionOfTransaction *mongo.Collection = database.OpenCollection(database.Client, os.Getenv("collection_name_transaction"))

func GatherTransactionRealTime(client *ethclient.Client, addresses []common.Address, chainID *big.Int) {

	var transactions []dataStruct.TransactionDetails

	for {

		// read block header
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		if LastBlockNumber == 0 || header.Number.Uint64() > LastBlockNumber {
			LastBlockNumber = header.Number.Uint64()

			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(LastBlockNumber)))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("block Number = ", LastBlockNumber)

			// count tx matched in block
			count := 0

			for _, tx := range block.Transactions() {

				signed := types.LatestSignerForChainID(chainID)
				from, err := types.Sender(signed, tx)
				if err != nil {
					log.Fatal(err)
					continue
				}

				// gathering all transaction match with list address
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
			if len(transactions) > 0 {
				storeTransaction(transactions)
			}

			//clear transaction array
			transactions = transactions[:0]

		}
		// read block each 12sec
		time.Sleep(12 * time.Second)

	}

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

func UnsubscribeAddress(address dataStruct.Address) error {
	filter := bson.M{"address": address.AddressSub}
	update := bson.M{"$set": bson.M{"status": false}}

	_, err := CollectionOfSubscriber.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
