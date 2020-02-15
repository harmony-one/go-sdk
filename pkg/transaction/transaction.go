package transaction

import (
	"math/big"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/numeric"
)

type Transaction = types.Transaction

func NewTransaction(
	nonce, gasLimit uint64,
	to address.T,
	shardID, toShardID uint32,
	amount, gasPrice numeric.Dec,
	data []byte) *Transaction {
	return types.NewCrossShardTransaction(nonce, &to, shardID, toShardID, amount.TruncateInt(), gasLimit, gasPrice.TruncateInt(), data[:])
}

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

func IsValid(tx *Transaction) bool {
	return true
}
