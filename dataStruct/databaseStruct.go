package dataStruct

type Address struct {
	AddressSub string `bson:"address"`
	Status     bool   `bson:"status"`
}
