package common

import (
	"fmt"
	"math/big"
)

//TODO Use reflection for these values instead of switch & slice given by AllChainIDs

// ChainID is a wrapper around the human name for a chain and the actual Big.Int used
type ChainID struct {
	Name  string
	Value *big.Int
}

type chainIDList struct {
	MainNet ChainID
	TestNet ChainID
}

// Chain is an enumeration of the known Chain-IDs
var Chain = chainIDList{
	MainNet: ChainID{"mainnet", big.NewInt(1)},
	TestNet: ChainID{"testnet", big.NewInt(2)},
}

// AllChainIDs returns list of known chains
func AllChainIDs() []string {
	return []string{"mainnet", "testnet"}
}

// StringToChainID returns the ChainID wrapper for the given human name of a chain-id
func StringToChainID(name string) (*ChainID, error) {
	switch name {
	case "mainnet":
		return &Chain.MainNet, nil
	case "testnet":
		return &Chain.TestNet, nil
	default:
		return nil, fmt.Errorf("unknown chain-id: %s", name)
	}
}
