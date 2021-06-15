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
	"github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/numeric"
)

var (
	nanoAsDec = numeric.NewDec(denominations.Nano)
	oneAsDec  = numeric.NewDec(denominations.One)

	// ErrBadTransactionParam is returned when invalid params are given to the
	// controller upon execution of a transaction.
	ErrBadTransactionParam = errors.New("transaction has bad parameters")
)

type p []interface{}

type transactionForRPC struct {
	params      map[string]interface{}
	transaction *types.Transaction
	// Hex encoded
	signature       *string
	transactionHash *string
	receipt         rpc.Reply
}

type sender struct {
	ks      *keystore.KeyStore
	account *accounts.Account
}

// Controller drives the transaction signing process
type Controller struct {
	executionError    error
	transactionErrors Errors
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
	handler rpc.T, senderKs *keystore.KeyStore,
	senderAcct *accounts.Account, chain common.ChainID,
	options ...func(*Controller),
) *Controller {
	txParams := make(map[string]interface{})
	ctrlr := &Controller{
		executionError: nil,
		messenger:      handler,
		sender: sender{
			ks:      senderKs,
			account: senderAcct,
		},
		transactionForRPC: transactionForRPC{
			params:          txParams,
			signature:       nil,
			transactionHash: nil,
			receipt:         nil,
		},
		chain:    chain,
		Behavior: behavior{false, Software, 0},
	}
	for _, option := range options {
		option(ctrlr)
	}
	return ctrlr
}

// TransactionToJSON dumps JSON representation
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

func (C *Controller) TransactionInfo() *types.Transaction {
	return C.transactionForRPC.transaction.Copy()
}

// TransactionHash - the tx hash
func (C *Controller) TransactionHash() *string {
	return C.transactionForRPC.transactionHash
}

// Receipt - the tx receipt
func (C *Controller) Receipt() rpc.Reply {
	return C.transactionForRPC.receipt
}

// TransactionErrors - tx errors
func (C *Controller) TransactionErrors() Errors {
	return C.transactionErrors
}

func (C *Controller) setShardIDs(fromShard, toShard uint32) {
	if C.executionError != nil {
		return
	}
	C.transactionForRPC.params["from-shard"] = fromShard
	C.transactionForRPC.params["to-shard"] = toShard
}

func (C *Controller) setIntrinsicGas(gasLimit uint64) {
	if C.executionError != nil {
		return
	}
	C.transactionForRPC.params["gas-limit"] = gasLimit
}

func (C *Controller) setGasPrice(gasPrice numeric.Dec) {
	if C.executionError != nil {
		return
	}
	if gasPrice.Sign() == -1 {
		C.executionError = ErrBadTransactionParam
		errorMsg := fmt.Sprintf(
			"can't set negative gas price: %d", gasPrice,
		)
		C.transactionErrors = append(C.transactionErrors, &Error{
			ErrMessage:           &errorMsg,
			TimestampOfRejection: time.Now().Unix(),
		})
		return
	}
	C.transactionForRPC.params["gas-price"] = gasPrice.Mul(nanoAsDec)
}

func (C *Controller) setAmount(amount numeric.Dec) {
	if C.executionError != nil {
		return
	}
	if amount.Sign() == -1 {
		C.executionError = ErrBadTransactionParam
		errorMsg := fmt.Sprintf(
			"can't set negative amount: %d", amount,
		)
		C.transactionErrors = append(C.transactionErrors, &Error{
			ErrMessage:           &errorMsg,
			TimestampOfRejection: time.Now().Unix(),
		})
		return
	}
	balanceRPCReply, err := C.messenger.SendRPC(
		rpc.Method.GetBalance,
		p{address.ToBech32(C.sender.account.Address), "latest"},
	)
	if err != nil {
		C.executionError = err
		return
	}
	currentBalance, _ := balanceRPCReply["result"].(string)
	bal, _ := new(big.Int).SetString(currentBalance[2:], 16)
	balance := numeric.NewDecFromBigInt(bal)
	gasAsDec := C.transactionForRPC.params["gas-price"].(numeric.Dec)
	gasAsDec = gasAsDec.Mul(numeric.NewDec(int64(C.transactionForRPC.params["gas-limit"].(uint64))))
	amountInAtto := amount.Mul(oneAsDec)
	total := amountInAtto.Add(gasAsDec)

	if total.GT(balance) {
		balanceInOne := balance.Quo(oneAsDec)
		C.executionError = ErrBadTransactionParam
		errorMsg := fmt.Sprintf(
			"insufficient balance of %s in shard %d for the requested transfer of %s",
			balanceInOne.String(), C.transactionForRPC.params["from-shard"].(uint32), amount.String(),
		)
		C.transactionErrors = append(C.transactionErrors, &Error{
			ErrMessage:           &errorMsg,
			TimestampOfRejection: time.Now().Unix(),
		})
		return
	}
	C.transactionForRPC.params["transfer-amount"] = amountInAtto
}

func (C *Controller) setReceiver(receiver *string) {
	if receiver != nil {
		addr := address.Parse(*receiver)
		C.transactionForRPC.params["receiver"] = &addr
	}
}

func (C *Controller) setNewTransactionWithDataAndGas(data []byte) {
	var addP *address.T
	if C.transactionForRPC.params["receiver"] != nil {
		addP = C.transactionForRPC.params["receiver"].(*address.T)
	}

	if C.executionError != nil {
		return
	}
	C.transactionForRPC.transaction = NewTransaction(
		C.transactionForRPC.params["nonce"].(uint64),
		C.transactionForRPC.params["gas-limit"].(uint64),
		addP,
		C.transactionForRPC.params["from-shard"].(uint32),
		C.transactionForRPC.params["to-shard"].(uint32),
		C.transactionForRPC.params["transfer-amount"].(numeric.Dec),
		C.transactionForRPC.params["gas-price"].(numeric.Dec),
		data,
	)
}

func (C *Controller) signAndPrepareTxEncodedForSending() {
	if C.executionError != nil {
		return
	}
	signedTransaction, err :=
		C.sender.ks.SignTx(*C.sender.account, C.transactionForRPC.transaction, C.chain.Value)
	if err != nil {
		C.executionError = err
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

func (C *Controller) hardwareSignAndPrepareTxEncodedForSending() {
	if C.executionError != nil {
		return
	}
	enc, signerAddr, err := ledger.SignTx(C.transactionForRPC.transaction, C.chain.Value)
	if err != nil {
		C.executionError = err
		return
	}
	if strings.Compare(signerAddr, address.ToBech32(C.sender.account.Address)) != 0 {
		C.executionError = ErrBadTransactionParam
		errorMsg := "signature verification failed : sender address doesn't match with ledger hardware addresss"
		C.transactionErrors = append(C.transactionErrors, &Error{
			ErrMessage:           &errorMsg,
			TimestampOfRejection: time.Now().Unix(),
		})
		return
	}
	hexSignature := hexutil.Encode(enc)
	C.transactionForRPC.signature = &hexSignature
}

func (C *Controller) sendSignedTx() {
	if C.executionError != nil || C.Behavior.DryRun {
		return
	}
	reply, err := C.messenger.SendRPC(rpc.Method.SendRawTransaction, p{C.transactionForRPC.signature})
	if err != nil {
		C.executionError = err
		return
	}
	r, _ := reply["result"].(string)
	C.transactionForRPC.transactionHash = &r
}

func (C *Controller) txConfirmation() {
	if C.executionError != nil || C.Behavior.DryRun {
		return
	}
	if C.Behavior.ConfirmationWaitTime > 0 {
		txHash := *C.TransactionHash()
		start := int(C.Behavior.ConfirmationWaitTime)
		for {
			r, _ := C.messenger.SendRPC(rpc.Method.GetTransactionReceipt, p{txHash})
			if r["result"] != nil {
				C.transactionForRPC.receipt = r
				return
			}
			transactionErrors, err := GetError(txHash, C.messenger)
			if err != nil {
				errMsg := fmt.Sprintf(err.Error())
				C.transactionErrors = append(C.transactionErrors, &Error{
					TxHashID:             &txHash,
					ErrMessage:           &errMsg,
					TimestampOfRejection: time.Now().Unix(),
				})
			}
			C.transactionErrors = append(C.transactionErrors, transactionErrors...)
			if len(transactionErrors) > 0 {
				C.executionError = fmt.Errorf("error found for transaction hash: %s", txHash)
				return
			}
			if start < 0 {
				C.executionError = fmt.Errorf("could not confirm transaction after %d seconds", C.Behavior.ConfirmationWaitTime)
				return
			}
			time.Sleep(time.Second)
			start--
		}
	}
}

// ExecuteTransaction is the single entrypoint to execute a plain transaction.
// Each step in transaction creation, execution probably includes a mutation
// Each becomes a no-op if executionError occurred in any previous step
func (C *Controller) ExecuteTransaction(
	nonce, gasLimit uint64,
	to *string,
	shardID, toShardID uint32,
	amount, gasPrice numeric.Dec,
	inputData []byte,
) error {
	// WARNING Order of execution matters
	C.setShardIDs(shardID, toShardID)
	C.setIntrinsicGas(gasLimit)
	C.setGasPrice(gasPrice)
	C.setAmount(amount)
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
	return C.executionError
}

func (C *Controller) SignTransaction(
	nonce, gasLimit uint64,
	to *string,
	shardID, toShardID uint32,
	amount, gasPrice numeric.Dec,
	inputData []byte,
) error {
	// WARNING Order of execution matters
	C.setShardIDs(shardID, toShardID)
	C.setIntrinsicGas(gasLimit)
	C.setGasPrice(gasPrice)
	C.setAmount(amount)
	C.setReceiver(to)
	C.transactionForRPC.params["nonce"] = nonce
	C.setNewTransactionWithDataAndGas(inputData)
	switch C.Behavior.SigningImpl {
	case Software:
		C.signAndPrepareTxEncodedForSending()
	case Ledger:
		C.hardwareSignAndPrepareTxEncodedForSending()
	}

	return C.executionError
}

// TODO: add logic to create staking transactions in the SDK.
