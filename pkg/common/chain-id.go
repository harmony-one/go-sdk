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
	MainNet    ChainID
	TestNet    ChainID
	PangaeaNet ChainID
	PartnerNet ChainID
	StressNet  ChainID
}

// Chain is an enumeration of the known Chain-IDs
var Chain = chainIDList{
	MainNet:    ChainID{"mainnet", big.NewInt(1)},
	TestNet:    ChainID{"testnet", big.NewInt(2)},
	PangaeaNet: ChainID{"pangaea", big.NewInt(3)},
	PartnerNet: ChainID{"partner", big.NewInt(4)},
	StressNet:  ChainID{"stress", big.NewInt(5)},
}

// AllChainIDs returns list of known chains
func AllChainIDs() []string {
	return []string{"mainnet", "testnet", "pangaea", "devnet", "partner", "stressnet"}
}

// StringToChainID returns the ChainID wrapper for the given human name of a chain-id
func StringToChainID(name string) (*ChainID, error) {
	switch name {
	case "mainnet":
		return &Chain.MainNet, nil
	case "testnet":
		return &Chain.TestNet, nil
	case "pangaea":
		return &Chain.PangaeaNet, nil
	case "devnet":
		return &Chain.PartnerNet, nil
	case "partner":
		return &Chain.PartnerNet, nil
	case "stressnet":
		return &Chain.StressNet, nil
	default:
		return nil, fmt.Errorf("unknown chain-id: %s", name)
	}
}
