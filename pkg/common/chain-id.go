package common

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// ChainID is a wrapper around the human name for a chain and the actual Big.Int used
type ChainID struct {
	Name  string   `json:"-"`
	Value *big.Int `json:"chain-as-number"`
}

type chainIDList struct {
	MainNet    ChainID `json:"mainnet"`
	TestNet    ChainID `json:"testnet"`
	PangaeaNet ChainID `json:"pangaea"`
	PartnerNet ChainID `json:"partner"`
	StressNet  ChainID `json:"stress"`
}

// Chain is an enumeration of the known Chain-IDs
var Chain = chainIDList{
	MainNet:    ChainID{"mainnet", big.NewInt(1)},
	TestNet:    ChainID{"testnet", big.NewInt(2)},
	PangaeaNet: ChainID{"pangaea", big.NewInt(3)},
	PartnerNet: ChainID{"partner", big.NewInt(4)},
	StressNet:  ChainID{"stress", big.NewInt(5)},
}

func (c chainIDList) String() string {
	s, _ := json.MarshalIndent(c, "", "  ")
	return string(s)
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
	case "dryrun":
		return &Chain.MainNet, nil
	default:
		return nil, fmt.Errorf("unknown chain-id: %s", name)
	}
}
