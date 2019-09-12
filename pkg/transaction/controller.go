package transaction

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/core"
)

type p []interface{}

type transactionForRPC struct {
	params      map[string]interface{}
	transaction *Transaction
	// Hex encoded
	signature *string
	receipt   *string
}

type sender struct {
	ks      *keystore.KeyStore
	account *accounts.Account
}

type Controller struct {
	failure           error
	messenger         rpc.T
	sender            sender
	transactionForRPC transactionForRPC
	chain             common.ChainID
	Behavior          behavior
}

type behavior struct {
	DryRun bool
}

func NewController(
	handler rpc.T,
	senderKs *keystore.KeyStore,
	senderAcct *accounts.Account,
	chain common.ChainID,
	options ...func(*Controller)) *Controller {

	txParams := make(map[string]interface{})
	ctrlr := &Controller{
		failure:   nil,
		messenger: handler,
		sender: sender{
			ks:      senderKs,
			account: senderAcct,
		},
		transactionForRPC: transactionForRPC{
			params:    txParams,
			signature: nil,
			receipt:   nil,
		},
		chain:    chain,
		Behavior: behavior{false},
	}
	for _, option := range options {
		option(ctrlr)
	}
	return ctrlr
}

func (C *Controller) verifyBalance(amount float64) {
	if C.failure != nil {
		return
	}
	balanceRPCReply, err := C.messenger.SendRPC(
		rpc.Method.GetBalance,
		p{address.ToBech32(C.sender.account.Address), "latest"},
	)
	if err != nil {
		C.failure = err
		return
	}
	currentBalance, _ := balanceRPCReply["result"].(string)
	balance, _ := big.NewInt(0).SetString(currentBalance[2:], 16)
	balance = common.NormalizeAmount(balance)
	transfer := big.NewInt(int64(amount * denominations.Nano))

	tns := (float64(transfer.Uint64()) / denominations.Nano)
	bln := (float64(balance.Uint64()) / denominations.Nano)

	if tns > bln {
		C.failure = errors.New(
			fmt.Sprintf("current balance of %.6f is not enough for the requested transfer %.6f", bln, tns),
		)
	}
}

func (C *Controller) setNextNonce() {
	if C.failure != nil {
		return
	}
	transactionCountRPCReply, err :=
		C.messenger.SendRPC(rpc.Method.GetTransactionCount, p{C.sender.account.Address.Hex(), "latest"})
	if err != nil {
		C.failure = err
		return
	}
	transactionCount, _ := transactionCountRPCReply["result"].(string)
	nonce, _ := big.NewInt(0).SetString(transactionCount[2:], 16)
	C.transactionForRPC.params["nonce"] = nonce.Uint64()
}

func (C *Controller) sendSignedTx() {
	if C.failure != nil || C.Behavior.DryRun {
		return
	}
	reply, err := C.messenger.SendRPC(rpc.Method.SendRawTransaction, p{C.transactionForRPC.signature})
	if err != nil {
		C.failure = err
		return
	}
	r, _ := reply["result"].(string)
	C.transactionForRPC.receipt = &r
}

func (C *Controller) setIntrinsicGas(rawInput string) {
	if C.failure != nil {
		return
	}
	inputData, _ := base64.StdEncoding.DecodeString(rawInput)
	gas, _ := core.IntrinsicGas(inputData, false, true)
	C.transactionForRPC.params["gas"] = gas
}

func (C *Controller) setGasPrice() {
	if C.failure != nil {
		return
	}
	C.transactionForRPC.params["gas-price"] = nil
}

func (C *Controller) setAmount(amount float64) {
	amountBigInt := big.NewInt(int64(amount * denominations.Nano))
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))
	C.transactionForRPC.params["transfer-amount"] = amt
}

func (C *Controller) setReceiver(receiver string) {
	C.transactionForRPC.params["receiver"] = address.Parse(receiver)
}

func (C *Controller) setNewTransactionWithData(i string, a float64) {
	if C.failure != nil {
		return
	}
	amountBigInt := big.NewInt(int64(a * denominations.Nano))
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))
	tx := NewTransaction(
		C.transactionForRPC.params["nonce"].(uint64),
		C.transactionForRPC.params["gas"].(uint64),
		C.transactionForRPC.params["receiver"].(address.T),
		C.transactionForRPC.params["from-shard"].(uint32),
		C.transactionForRPC.params["to-shard"].(uint32),
		amt,
		C.chain.Value,
		// C.transactionForRPC.params["gas-price"].(*big.Int),
		[]byte(i),
	)
	if common.DebugTransaction {
		r, _ := tx.MarshalJSON()
		fmt.Println(string(r))
	}
	C.transactionForRPC.transaction = tx
}

func (C *Controller) signAndPrepareTxEncodedForSending() {
	if C.failure != nil {
		return
	}
	signedTransaction, err :=
		C.sender.ks.SignTx(*C.sender.account, C.transactionForRPC.transaction, C.chain.Value)
	if err != nil {
		fmt.Println(err)
	}
	enc, _ := rlp.EncodeToBytes(signedTransaction)
	hexSignature := hexutil.Encode(enc)
	C.transactionForRPC.signature = &hexSignature
}

func (C *Controller) setShardIDs(fromShard, toShard int) {
	if C.failure != nil {
		return
	}
	C.transactionForRPC.params["from-shard"] = uint32(fromShard)
	C.transactionForRPC.params["to-shard"] = uint32(toShard)
}

func (C *Controller) Receipt() *string {
	return C.transactionForRPC.receipt
}

func (C *Controller) ExecuteTransaction(
	to, inputData string,
	amount float64,
	fromShard, toShard int,
) error {
	C.transactionForRPC.params["gas-price"] = big.NewInt(0)
	// WARNING Order of execution matters
	C.setShardIDs(fromShard, toShard)
	C.setIntrinsicGas(inputData)
	C.setAmount(amount)
	C.verifyBalance(amount)
	C.setReceiver(to)
	C.setGasPrice()
	C.setNextNonce()
	C.setNewTransactionWithData(inputData, amount)
	C.signAndPrepareTxEncodedForSending()
	C.sendSignedTx()
	return C.failure
}
