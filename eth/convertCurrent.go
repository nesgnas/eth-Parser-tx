package eth

import (
	"math"
	"math/big"
)

func ConvertWei(valueETH *big.Int) string {
	concurrent := new(big.Float)
	concurrent.SetString(valueETH.String())
	ethValue := new(big.Float).Quo(concurrent, big.NewFloat(math.Pow10(18)))
	return ethValue.String()
}
