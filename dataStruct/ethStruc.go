package dataStruct

type TransactionDetails struct {
	BlockNumber uint64 `bson:"block_number"`
	TxHash      string `bson:"hash"`
	From        string `bson:"from"`
	To          string `bson:"to"`
	Value       string `bson:"value"`
}

var TrasactionGlobal []TransactionDetails
