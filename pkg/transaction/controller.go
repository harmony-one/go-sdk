package transaction

import (
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
	"github.com/harmony-one/harmony/numeric"
)

var (
	nanoAsDec = numeric.NewDec(denominations.Nano)
	oneAsDec  = numeric.NewDec(denominations.One)
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

func (C *Controller) verifyBalance() {
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
	bal, _ := new(big.Int).SetString(currentBalance[2:], 16)
	balance := numeric.NewDecFromBigInt(bal)
	gasAsDec := C.transactionForRPC.params["gas-price"].(numeric.Dec)
	gasAsDec = gasAsDec.Mul(numeric.NewDec(int64(C.transactionForRPC.params["gas-limit"].(uint64))))
	total := C.transactionForRPC.params["transfer-amount"].(numeric.Dec).Add(gasAsDec)

	if total.GT(balance) {
		b := balance.Quo(oneAsDec)
		t := total.Quo(oneAsDec)
		C.failure = fmt.Errorf(
			"insufficient balance of %s in shard %d for the requested transfer of %s", b.String(), C.transactionForRPC.params["from-shard"].(uint32), t.String(),
		)
	}
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

func (C *Controller) setIntrinsicGas(gasLimit uint64) {
	if C.failure != nil {
		return
	}
	C.transactionForRPC.params["gas-limit"] = gasLimit
}

func (C *Controller) setGasPrice(gasPrice numeric.Dec) {
	if C.failure != nil {
		return
	}
	C.transactionForRPC.params["gas-price"] = gasPrice.Mul(nanoAsDec)
}

func (C *Controller) setAmount(amount numeric.Dec) {
	amt := amount.Mul(oneAsDec)
	C.transactionForRPC.params["transfer-amount"] = amt
}

func (C *Controller) setReceiver(receiver string) {
	C.transactionForRPC.params["receiver"] = address.Parse(receiver)
}

func (C *Controller) setNewTransactionWithDataAndGas(i string) {
	if C.failure != nil {
		return
	}

	tx := NewTransaction(
		C.transactionForRPC.params["nonce"].(uint64),
		C.transactionForRPC.params["gas-limit"].(uint64),
		C.transactionForRPC.params["receiver"].(address.T),
		C.transactionForRPC.params["from-shard"].(uint32),
		C.transactionForRPC.params["to-shard"].(uint32),
		C.transactionForRPC.params["transfer-amount"].(numeric.Dec),
		C.transactionForRPC.params["gas-price"].(numeric.Dec),
		[]byte(i),
	)
	C.transactionForRPC.transaction = tx
}

// TransactionToJSON dumps JSON rep
func (C *Controller) TransactionToJSON(pretty bool) string {
	r, _ := C.transactionForRPC.transaction.MarshalJSON()
	if pretty {
		return common.JSONPrettyFormat(string(r))
	}
	return string(r)
}

// RawTransaction dumps the signature as string
func (C *Controller) RawTransaction() string {
	return *C.transactionForRPC.signature
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
		fmt.Println("Signed with ChainID:", C.transactionForRPC.transaction.ChainID())
		fmt.Println(common.JSONPrettyFormat(string(r)))
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
	if C.failure != nil || C.Behavior.DryRun {
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

// ExecuteTransaction is the single entrypoint to execute a transaction.
// Each step in transaction creation, execution probably includes a mutation
// Each becomes a no-op if failure occured in any previous step
func (C *Controller) ExecuteTransaction(
	to, inputData string,
	amount, gasPrice numeric.Dec, nonce, gasLimit uint64,
	fromShard, toShard int,
) error {
	// WARNING Order of execution matters
	C.setShardIDs(fromShard, toShard)
	C.setIntrinsicGas(gasLimit)
	C.setGasPrice(gasPrice)
	C.setAmount(amount)
	C.verifyBalance()
	C.setReceiver(to)
	C.transactionForRPC.params["nonce"] = nonce
	C.setNewTransactionWithDataAndGas(inputData)
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
