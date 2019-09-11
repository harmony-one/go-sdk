package common

import (
	"math/big"
	"os"

	"github.com/harmony-one/harmony/accounts/keystore"
)

const (
	DefaultConfigDirName               = ".hmy_cli"
	DefaultConfigAccountAliasesDirName = "account-keys"
	DefaultPassphrase                  = "harmony-one"
	JSONRPCVersion                     = "2.0"
)

var (
	ScryptN          = keystore.StandardScryptN
	ScryptP          = keystore.StandardScryptP
	DebugRPC         = false
	DebugTransaction = false
)

func init() {
	if _, enabled := os.LookupEnv("HMY_RPC_DEBUG"); enabled != false {
		DebugRPC = true
	}
	if _, enabled := os.LookupEnv("HMY_TX_DEBUG"); enabled != false {
		DebugTransaction = true
	}

}

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
