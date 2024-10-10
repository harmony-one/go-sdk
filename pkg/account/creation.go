package account

import (
	"errors"

	"github.com/harmony-one/go-sdk/pkg/keys"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	"github.com/harmony-one/go-sdk/pkg/store"
)

var (
	AccountByNameExists = errors.New("name chosen for account already exists")
)

type Creation struct {
	Name            string
	Passphrase      string
	Mnemonic        string
	HdAccountNumber *uint32
	HdIndexNumber   *uint32
	CoinType        *uint32
}

func New() string {
	return "New Account"
}

func IsValidPassphrase(pass string) bool {
	return true
}

// CreateNewLocalAccount assumes all the inputs are valid, legitmate
func CreateNewLocalAccount(candidate *Creation) error {
	ks := store.FromAccountName(candidate.Name)
	if candidate.Mnemonic == "" {
		candidate.Mnemonic = mnemonic.Generate()
	}

	index := uint32(0)
	if candidate.HdIndexNumber != nil {
		index = *candidate.HdIndexNumber
	}

	coinType := uint32(1023)
	if candidate.CoinType != nil {
		coinType = *candidate.CoinType
	}

	private, _ := keys.FromMnemonicSeedAndPassphrase(candidate.Mnemonic, int(index), int(coinType))
	_, err := ks.ImportECDSA(private.ToECDSA(), candidate.Passphrase)
	if err != nil {
		return err
	}
	return nil
}
