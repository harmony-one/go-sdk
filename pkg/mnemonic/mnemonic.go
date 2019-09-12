package mnemonic

import (
	"errors"

	"github.com/tyler-smith/go-bip39"
)

var (
	InvalidMnemonic = errors.New("invalid mnemonic given")
)

func Generate() string {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}
