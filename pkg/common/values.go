package common

import (
	"math/big"

	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/go-sdk/pkg/address"
)

const (
	DefaultConfigDirName               = ".hmy_cli"
	DefaultConfigAccountAliasesDirName = "account-keys"
	DefaultPassphrase                  = "harmony-one"
)

var (
	ScryptN = keystore.StandardScryptN
	ScryptP = keystore.StandardScryptP
)

type ChainID struct {
	Name  string
	Value *big.Int
}

type chainIDList struct {
	MainNet ChainID
	TestNet ChainID
}

var Chain = chainIDList{
	MainNet: ChainID{"mainnet", big.NewInt(1)},
	TestNet: ChainID{"testnet", big.NewInt(2)},
}

func StringToChainID(name string) *ChainID {
	switch name {
	case "mainnet":
		return &Chain.MainNet
	case "testnet":
		return &Chain.TestNet
	default:
		return nil
	}
}

/////////// OneAddress

// OneAddress type for validates address entered with cli_flags
type OneAddress string

func (oneAddress OneAddress) String() string {
	return string(oneAddress)
}
// Set and validate OneAddress
func (oneAddress OneAddress) Set(s string) error {
	_, err := address.Bech32ToAddress(s)
	if err != nil {
		return err;
	}
	return nil;
}
// Type of OneAddress
func (oneAddress OneAddress) Type() string {
	return "OneAddress"
}