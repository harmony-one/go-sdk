package transaction

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/ledger"
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
	signature   *string
	receiptHash *string
	receipt     rpc.Reply
}

type sender struct {
	ks      *keystore.KeyStore
	account *accounts.Account
}

// Controller drives the transaction signing process
type Controller struct {
	failure           error
	messenger         rpc.T
	sender            sender
	transactionForRPC transactionForRPC
	chain             common.ChainID
	Behavior          behavior
}

type behavior struct {
	DryRun               bool
	SigningImpl          SignerImpl
	ConfirmationWaitTime uint32
}

// NewController initializes a Controller, caller can control behavior via options
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
			params:      txParams,
			signature:   nil,
			receiptHash: nil,
			receipt:     nil,
		},
		chain:    chain,
		Behavior: behavior{false, Software, 0},
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
		C.failure = fmt.Errorf(
			"current balance of %.6f is not enough for the requested transfer %.6f", bln, tns,
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
	C.transactionForRPC.receiptHash = &r
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

func (C *Controller) setNewTransactionWithDataAndGas(i string, amount, gasPrice float64) {
	if C.failure != nil {
		return
	}
	amountBigInt := big.NewInt(int64(amount * denominations.Nano))
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))
	gPrice := big.NewInt(int64(gasPrice))
	gPrice = gPrice.Mul(gPrice, big.NewInt(denominations.Nano))

	tx := NewTransaction(
		C.transactionForRPC.params["nonce"].(uint64),
		C.transactionForRPC.params["gas"].(uint64),
		C.transactionForRPC.params["receiver"].(address.T),
		C.transactionForRPC.params["from-shard"].(uint32),
		C.transactionForRPC.params["to-shard"].(uint32),
		amt,
		gPrice,
		[]byte(i),
	)
	C.transactionForRPC.transaction = tx
}

func (C *Controller) TransactionToJSON(pretty bool) string {
	r, _ := C.transactionForRPC.transaction.MarshalJSON()
	if pretty {
		return common.JSONPrettyFormat(string(r))
	}
	return string(r)
}

func (C *Controller) signAndPrepareTxEncodedForSending() {
	if C.failure != nil {
		return
	}
	signedTransaction, err :=
		C.sender.ks.SignTx(*C.sender.account, C.transactionForRPC.transaction, C.chain.Value)
	if err != nil {
		C.failure = err
		return
	}
	C.transactionForRPC.transaction = signedTransaction
	enc, _ := rlp.EncodeToBytes(signedTransaction)
	hexSignature := hexutil.Encode(enc)
	C.transactionForRPC.signature = &hexSignature
	if common.DebugTransaction {
		r, _ := signedTransaction.MarshalJSON()
		fmt.Println(string(r))
	}
}

func (C *Controller) setShardIDs(fromShard, toShard int) {
	if C.failure != nil {
		return
	}
	C.transactionForRPC.params["from-shard"] = uint32(fromShard)
	C.transactionForRPC.params["to-shard"] = uint32(toShard)
}

func (C *Controller) ReceiptHash() *string {
	return C.transactionForRPC.receiptHash
}

func (C *Controller) Receipt() rpc.Reply {
	return C.transactionForRPC.receipt
}

func (C *Controller) hardwareSignAndPrepareTxEncodedForSending() {
	if C.failure != nil {
		return
	}
	enc, signerAddr, err := ledger.SignTx(C.transactionForRPC.transaction, C.chain.Value)
	if err != nil {
		C.failure = err
		return
	}
	if strings.Compare(signerAddr, address.ToBech32(C.sender.account.Address)) != 0 {
		C.failure = errors.New("signature verification failed : sender address doesn't match with ledger hardware addresss")
		return
	}
	hexSignature := hexutil.Encode(enc)
	C.transactionForRPC.signature = &hexSignature
}

func (C *Controller) txConfirmation() {
	if C.failure != nil {
		return
	}
	if C.Behavior.ConfirmationWaitTime > 0 {
		receipt := *C.ReceiptHash()
		start := int(C.Behavior.ConfirmationWaitTime)
		for {
			if start < 0 {
				return
			}
			r, _ := C.messenger.SendRPC(rpc.Method.GetTransactionReceipt, p{receipt})
			if r["result"] != nil {
				C.transactionForRPC.receipt = r
				return
			}
			time.Sleep(time.Second * 2)
			start = start - 2
		}
	}
}

func (C *Controller) ExecuteTransaction(
	to, inputData string,
	amount, gPrice float64,
	fromShard, toShard int,
) error {
	// WARNING Order of execution matters
	C.setShardIDs(fromShard, toShard)
	C.setIntrinsicGas(inputData)
	C.setAmount(amount)
	C.verifyBalance(amount)
	C.setReceiver(to)
	C.setGasPrice()
	C.setNextNonce()
	C.setNewTransactionWithDataAndGas(inputData, amount, gPrice)
	switch C.Behavior.SigningImpl {
	case Software:
		C.signAndPrepareTxEncodedForSending()
	case Ledger:
		C.hardwareSignAndPrepareTxEncodedForSending()
	}
	C.sendSignedTx()
	C.txConfirmation()
	return C.failure
}
