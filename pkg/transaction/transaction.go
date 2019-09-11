package transaction

import (
	"math/big"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/harmony/core/types"
)

type Transaction = types.Transaction

func NewTransaction(
	nonce, gasLimit uint64,
	to address.T,
	shardID, toShardID uint32,
	amount, gasPrice *big.Int,
	data []byte) *Transaction {
	// types.New
	return types.NewCrossShardTransaction(nonce, &to, shardID, toShardID, amount, gasLimit, gasPrice, data)
}

func IsValid(tx *Transaction) bool {
	return true
}
