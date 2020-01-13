package common

import (
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
)

var (
	pattern, _ = regexp.Compile("[0-9]+\\.{0,1}[0-9]*e-{0,1}[0-9]+")
)

func NormalizeAmount(value *big.Int) *big.Int {
	return value.Div(value, big.NewInt(denominations.Nano))
}

func Pow(base numeric.Dec, exp int) numeric.Dec {
	if (exp < 0) {
		return Pow(numeric.NewDec(1).Quo(base), -exp)
	}
	result := numeric.NewDec(1)
	for {
		if (exp % 2 == 1) {
			result = result.Mul(base)
		}
		exp = exp >> 1
		if (exp == 0) {
			break
		}
		base = base.Mul(base)
	}
	return result
}

func NewDecFromString (i string) (numeric.Dec, error) {
	if pattern.FindString(i) != "" {
		tokens := strings.Split(i, "e")
		a, _ := numeric.NewDecFromStr(tokens[0])
		b, _ := strconv.Atoi(tokens[1])
		return a.Mul(Pow(numeric.NewDec(10), b)), nil
	} else {
		return numeric.NewDecFromStr(i)
	}
}
