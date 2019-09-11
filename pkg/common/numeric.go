package common

import (
	"math/big"

	"github.com/harmony-one/harmony/common/denominations"
)

func NormalizeAmount(value *big.Int) *big.Int {
	return value.Div(value, big.NewInt(denominations.Nano))
}
