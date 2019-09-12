package account

import (
	"errors"
	"fmt"

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

const (
	TestMnemonic   = "quick parade stay hockey build token access sentence choice supply creek twelve"
	TestPassphrase = "edgar"
	TestPrivateKey = "0xc83827745d4e2e74e0b996b2a31ebff7bde72790bf23db839668223c07fb0299"
	TestPublicKey  = "0x038044dbdd8e6f901d26e3e705570b87fa6ca412f435979c60e90539a5b72b5a9f"
)

// By this point assume all the inputs are valid, legitmate
func CreateNewLocalAccount(candidate *Creation) error {
	ks := store.FromAccountName(candidate.Name)
	if candidate.Mnemonic == "" {
		candidate.Mnemonic = mnemonic.Generate()
	}
	private, public := keys.FromMnemonicSeedAndPassphrase(candidate.Mnemonic, candidate.Passphrase)
	acct, err := ks.ImportECDSA(private.ToECDSA(), candidate.Passphrase)
	if err != nil {
		fmt.Println(acct.Address.Hex(), public)
		return err
	}
	return nil
}
