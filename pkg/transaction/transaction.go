package transaction

import (
	"math/big"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/numeric"
)

// NewTransaction - create a new Transaction based on supplied params
func NewTransaction(
	nonce, gasLimit uint64,
	to *address.T,
	shardID, toShardID uint32,
	amount, gasPrice numeric.Dec,
	data []byte) *types.Transaction {
	return types.NewCrossShardTransaction(nonce, to, shardID, toShardID, amount.TruncateInt(), gasLimit, gasPrice.TruncateInt(), data[:])
}

// NewEthTransaction - create a new Transaction based on supplied params
func NewEthTransaction(
	nonce, gasLimit uint64,
	to address.T,
	amount, gasPrice numeric.Dec,
	data []byte) *types.EthTransaction {
	return types.NewEthTransaction(nonce, to, amount.TruncateInt(), gasLimit, gasPrice.TruncateInt(), data[:])
}

// GetNextNonce returns the nonce on-chain (finalized transactions)
func GetNextNonce(addr string, messenger rpc.T) uint64 {
	transactionCountRPCReply, err :=
		messenger.SendRPC(rpc.Method.GetTransactionCount, []interface{}{address.Parse(addr), "latest"})

	if err != nil {
		return 0
	}

	transactionCount, _ := transactionCountRPCReply["result"].(string)
	n, _ := big.NewInt(0).SetString(transactionCount[2:], 16)
	return n.Uint64()
}

// GetNextPendingNonce returns the nonce from the tx-pool (un-finalized transactions)
func GetNextPendingNonce(addr string, messenger rpc.T) uint64 {
	transactionCountRPCReply, err :=
		messenger.SendRPC(rpc.Method.GetTransactionCount, []interface{}{address.Parse(addr), "pending"})

	if err != nil {
		return 0
	}

	transactionCount, _ := transactionCountRPCReply["result"].(string)
	n, _ := big.NewInt(0).SetString(transactionCount[2:], 16)
	return n.Uint64()
}

// IsValid - whether or not a tx is valid
func IsValid(tx *types.Transaction) bool {
	return true
}
