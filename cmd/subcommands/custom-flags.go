package cmd

import (
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/validation"
	"github.com/pkg/errors"
)

type oneAddress struct {
	address string
}

func (oneAddress oneAddress) String() string {
	return oneAddress.address
}

func (oneAddress *oneAddress) Set(s string) error {
	err := validation.ValidateAddress(s)
	if err != nil {
		return err
	}

	_, err = address.Bech32ToAddress(s)
	if err != nil {
		return errors.Wrap(err, "not a valid one address")
	}
	oneAddress.address = s
	return nil
}

func (oneAddress oneAddress) Type() string {
	return "string"
}

type chainIDWrapper struct {
	chainID *common.ChainID
}

func (chainIDWrapper chainIDWrapper) String() string {
	return chainIDWrapper.chainID.Name
}

func (chainIDWrapper *chainIDWrapper) Set(s string) error {
	chain, err := common.StringToChainID(s)
	chainIDWrapper.chainID = chain
	if err != nil {
		return err
	}
	return nil
}

func (chainIDWrapper chainIDWrapper) Type() string {
	return "string"
}
