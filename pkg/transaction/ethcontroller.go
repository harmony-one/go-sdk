package transaction

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/numeric"
)

type ethTransactionForRPC struct {
	params      map[string]interface{}
	transaction *types.EthTransaction
	// Hex encoded
	signature       *string
	transactionHash *string
	receipt         rpc.Reply
}

// EthController drives the eth transaction signing process
type EthController struct {
	executionError    error
	transactionErrors Errors
	messenger         rpc.T
	sender            sender
	transactionForRPC ethTransactionForRPC
	chain             common.ChainID
	Behavior          behavior
}

// NewEthController initializes a EthController, caller can control behavior via options
func NewEthController(
	handler rpc.T, senderKs *keystore.KeyStore,
	senderAcct *accounts.Account, chain common.ChainID,
	options ...func(*EthController),
) *EthController {
	txParams := make(map[string]interface{})
	ctrlr := &EthController{
		executionError: nil,
		messenger:      handler,
		sender: sender{
			ks:      senderKs,
			account: senderAcct,
		},
		transactionForRPC: ethTransactionForRPC{
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

// EthTransactionToJSON dumps JSON representation
func (C *EthController) EthTransactionToJSON(pretty bool) string {
	r, _ := C.transactionForRPC.transaction.MarshalJSON()
	if pretty {
		return common.JSONPrettyFormat(string(r))
	}
	return string(r)
}

// RawTransaction dumps the signature as string
func (C *EthController) RawTransaction() string {
	return *C.transactionForRPC.signature
}

// TransactionHash - the tx hash
func (C *EthController) TransactionHash() *string {
	return C.transactionForRPC.transactionHash
}

// Receipt - the tx receipt
func (C *EthController) Receipt() rpc.Reply {
	return C.transactionForRPC.receipt
}

// TransactionErrors - tx errors
func (C *EthController) TransactionErrors() Errors {
	return C.transactionErrors
}

func (C *EthController) setIntrinsicGas(gasLimit uint64) {
	if C.executionError != nil {
		return
	}
	C.transactionForRPC.params["gas-limit"] = gasLimit
}

func (C *EthController) setGasPrice(gasPrice numeric.Dec) {
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

func (C *EthController) setAmount(amount numeric.Dec) {
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

func (C *EthController) setReceiver(receiver string) {
	C.transactionForRPC.params["receiver"] = address.Parse(receiver)
}

func (C *EthController) setNewTransactionWithDataAndGas(data []byte) {
	if C.executionError != nil {
		return
	}
	C.transactionForRPC.transaction = NewEthTransaction(
		C.transactionForRPC.params["nonce"].(uint64),
		C.transactionForRPC.params["gas-limit"].(uint64),
		C.transactionForRPC.params["receiver"].(address.T),
		C.transactionForRPC.params["transfer-amount"].(numeric.Dec),
		C.transactionForRPC.params["gas-price"].(numeric.Dec),
		data,
	)
}

func (C *EthController) signAndPrepareTxEncodedForSending() {
	if C.executionError != nil {
		return
	}
	signedTransaction, err :=
		C.sender.ks.SignEthTx(*C.sender.account, C.transactionForRPC.transaction, C.chain.Value)
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

/*func (C *EthController) hardwareSignAndPrepareTxEncodedForSending() {
	if C.executionError != nil {
		return
	}
	enc, signerAddr, err := ledger.SignEthTx(C.transactionForRPC.transaction, C.chain.Value)
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
}*/

func (C *EthController) sendSignedTx() {
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

func (C *EthController) txConfirmation() {
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

// ExecuteEthTransaction is the single entrypoint to execute an eth transaction.
// Each step in transaction creation, execution probably includes a mutation
// Each becomes a no-op if executionError occurred in any previous step
func (C *EthController) ExecuteEthTransaction(
	nonce, gasLimit uint64,
	to string,
	amount, gasPrice numeric.Dec,
	inputData []byte,
) error {
	// WARNING Order of execution matters
	C.setIntrinsicGas(gasLimit)
	C.setGasPrice(gasPrice)
	C.setAmount(amount)
	C.setReceiver(to)
	C.transactionForRPC.params["nonce"] = nonce
	C.setNewTransactionWithDataAndGas(inputData)
	switch C.Behavior.SigningImpl {
	case Software:
		C.signAndPrepareTxEncodedForSending()
		/*case Ledger:
		C.hardwareSignAndPrepareTxEncodedForSending()*/
	}
	C.sendSignedTx()
	C.txConfirmation()
	return C.executionError
}

// TODO: add logic to create staking transactions in the SDK.
