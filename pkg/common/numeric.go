package common

import (
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/harmony-one/harmony/numeric"
)

var (
	pattern, _ = regexp.Compile("[0-9]+\\.{0,1}[0-9]*e-{0,1}[0-9]+")
)

func Pow(base numeric.Dec, exp int) numeric.Dec {
	if exp < 0 {
		return Pow(numeric.NewDec(1).Quo(base), -exp)
	}
	result := numeric.NewDec(1)
	for {
		if exp%2 == 1 {
			result = result.Mul(base)
		}
		exp = exp >> 1
		if exp == 0 {
			break
		}
		base = base.Mul(base)
	}
	return result
}

func NewDecFromString(i string) (numeric.Dec, error) {
	if pattern.FindString(i) != "" {
		tokens := strings.Split(i, "e")
		a, _ := numeric.NewDecFromStr(tokens[0])
		b, _ := strconv.Atoi(tokens[1])
		return a.Mul(Pow(numeric.NewDec(10), b)), nil
	} else {
		if strings.HasPrefix(i, ".") {
			i = "0" + i
		}
		return numeric.NewDecFromStr(i)
	}
}

// Assumes Hex string input
// Split into 2 64 bit integers to guarentee 128 bit precision
func NewDecFromHex(str string) numeric.Dec {
	str = strings.TrimPrefix(str, "0x")
	half := len(str) / 2
	right := str[half:]
	r, _ := big.NewInt(0).SetString(right, 16)
	if half == 0 {
		return numeric.NewDecFromBigInt(r)
	}
	left := str[:half]
	l, _ := big.NewInt(0).SetString(left, 16)
	return numeric.NewDecFromBigInt(l).Mul(
		Pow(numeric.NewDec(16), len(right)),
	).Add(numeric.NewDecFromBigInt(r))
}
