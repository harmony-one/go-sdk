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
}

func New() string {
	return "New Account"
}

func IsValidPassphrase(pass string) bool {
	return true
}

// By this point assume all the inputs are valid, legitmate
func CreateNewLocalAccount(candidate *Creation) error {
	ks := store.FromAccountName(candidate.Name)
	if candidate.Mnemonic == "" {
		candidate.Mnemonic = mnemonic.Generate()
	}
	private, _ := keys.FromMnemonicSeedAndPassphrase(candidate.Mnemonic, candidate.Passphrase)
	_, err := ks.ImportECDSA(private.ToECDSA(), candidate.Passphrase)
	if err != nil {
		return err
	}
	return nil
}
