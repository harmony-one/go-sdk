package common

import (
	"errors"
	"math/big"
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

var (
	unknownChain = errors.New("unknown chain id provided")
)

func StringToChainID(name string) (*ChainID, error) {
	switch name {
	case "mainnet":
		return &Chain.MainNet, nil
	case "testnet":
		return &Chain.TestNet, nil
	default:
		return nil, unknownChain
	}
}
