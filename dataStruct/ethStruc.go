package dataStruct

type TransactionDetails struct {
	BlockNumber uint64
	TxHash      string
	From        string
	To          string
	Value       string
}

var TrasactionGlobal []TransactionDetails
