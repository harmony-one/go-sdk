package common

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/common/address"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/core"
	"github.com/harmony-one/harmony/core/types"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	debugEnabled  = false
	defaultKeyDir string
)

func init() {
	if _, enabled := os.LookupEnv("HMY_TX_DEBUG"); enabled != false {
		debugEnabled = true
	}
	userDir, _ := homedir.Dir()
	defaultKeyDir = path.Join(userDir, ".hmy_cli", "keystore")
}

type sender struct {
	addr     common.Address
	account  accounts.Account
	txParams map[string]interface{}
	// nil means not ready for usage
	readyTransaction            *types.Transaction
	signedAndEncodedTransaction string
	txReceipt                   *uint64
}

type TxController struct {
	failure          error
	messenger        rpc.T
	PreferOneAddress bool
	WaitForTxConfirm bool
	ks               *keystore.KeyStore
	sender           sender
	receipt          *uint64
}

// TODO Make node more robust with URL validation
//TODO useOneAddressInsteadOfHex, waitTxConfirm bool as functional API params
func NewTxController(handler rpc.T, senderAddr string, options ...func(*TxController)) (*TxController, error) {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore("/Users/edgar/.hmy_cli/keystore", scryptN, scryptP)
	senderParsed := address.Parse(senderAddr)
	account, lookupError := ks.Find(accounts.Account{Address: senderParsed})
	unlockError := ks.Unlock(account, "")

	if lookupError != nil || unlockError != nil {
		return nil, errors.New("Lookup or account unlocking of sender address in local keystore failed")
	}

	ctrlr := &TxController{
		failure:   nil,
		messenger: handler,
		ks:        ks,
		sender:    sender{senderParsed, account, make(map[string]interface{}), nil, "", nil},
	}
	for _, option := range options {
		option(ctrlr)
	}

	return ctrlr, nil
}

func (C *TxController) verifyBalance(amount float64) {
	if C.failure != nil {
		return
	}
	balanceRPCReply := C.messenger.SendRPC(
		rpc.Method.GetBalance,
		[]interface{}{address.ToBech32(C.sender.addr), "latest"},
	)
	currentBalance, _ := balanceRPCReply["result"].(string)
	balance, _ := big.NewInt(0).SetString(currentBalance[2:], 16)
	balance = NormalizeAmount(balance)
	transfer := big.NewInt(int64(amount * denominations.Nano))

	tns := (float64(transfer.Uint64()) / denominations.Nano)
	bln := (float64(balance.Uint64()) / denominations.Nano)

	if tns > bln {
		C.failure = errors.New(
			fmt.Sprintf("current balance of %.6f is not enough for the requested transfer %.6f", bln, tns),
		)
	}
}

func (C *TxController) setNextNonce() {
	if C.failure != nil {
		return
	}
	transactionCountRPCReply := C.messenger.SendRPC(
		rpc.Method.GetTransactionCount,
		[]interface{}{C.sender.addr.Hex(), "latest"},
	)
	transactionCount, _ := transactionCountRPCReply["result"].(string)
	nonce, _ := big.NewInt(0).SetString(transactionCount[2:], 16)
	C.sender.txParams["nonce"] = nonce.Uint64()
}

func (C *TxController) sendSignedTx() {
	if C.failure != nil {
		return
	}
	reply := C.messenger.SendRPC(
		rpc.Method.SendRawTransaction,
		[]interface{}{C.sender.signedAndEncodedTransaction},
	)
	txReceipt, _ := reply["result"].(string)
	receipt, _ := big.NewInt(0).SetString(txReceipt[2:], 16)
	receiptAsUint := receipt.Uint64()
	C.sender.txReceipt = &receiptAsUint
}

func (C *TxController) setIntrinsicGas(rawInput string) {
	if C.failure != nil {
		return
	}
	inputData, _ := base64.StdEncoding.DecodeString(rawInput)
	gas, _ := core.IntrinsicGas(inputData, false, true)
	C.sender.txParams["gas"] = gas
}

func (C *TxController) setGasPrice() {
	if C.failure != nil {
		return
	}
	C.sender.txParams["gas-price"] = nil
}

func (C *TxController) setAmount(amount float64) {
	amountBigInt := big.NewInt(int64(amount * denominations.Nano))
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))
	C.sender.txParams["transfer-amount"] = amt
}

func (C *TxController) setReceiver(receiver string) {
	C.sender.txParams["receiver"] = address.Parse(receiver)
}

func (C *TxController) setNewTransactionWithData(inputData string, amount float64) {
	if C.failure != nil {
		return
	}
	// TODO Refactor to use the cross-shard transaction item
	// tx := hmyTypes.NewCrossShardTransaction(
	// 	transactionCount, &receiverAddress, fromShard, toShard, amountBigInt,
	// 	gas, gasPriceBigInt, inputData)

	var gasPrice *big.Int = nil
	if value, ok := C.sender.txParams["gas-price"].(*big.Int); ok {
		gasPrice = value
	}

	amountBigInt := big.NewInt(int64(amount * denominations.Nano))
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

	tx := types.NewTransaction(
		C.sender.txParams["nonce"].(uint64),
		C.sender.txParams["receiver"].(common.Address),
		C.sender.txParams["shardID"].(uint32),
		amt,
		// C.sender.txParams["transfer-amount"].(*big.Int),
		C.sender.txParams["gas"].(uint64),
		gasPrice,
		[]byte(inputData),
	)
	if debugEnabled {
		r, _ := tx.MarshalJSON()
		fmt.Println(string(r))
	}
	C.sender.readyTransaction = tx
}

func (C *TxController) signAndPrepareTxEncodedForSending() {
	if C.failure != nil {
		return
	}

	signedTransaction, err := C.ks.SignTx(C.sender.account, C.sender.readyTransaction, nil)
	if err != nil {
		fmt.Println(err)
	}
	enc, _ := rlp.EncodeToBytes(signedTransaction)
	rawTx := hexutil.Encode(enc)
	C.sender.signedAndEncodedTransaction = rawTx
}

func (C *TxController) ExecuteTransaction(to, inputData string, amount float64, fromShard, toShard int) error {
	fmt.Println(to, inputData, amount, fromShard, toShard)
	// HACK This is pre cross shard transaction
	C.sender.txParams["shardID"] = uint32(0)
	// WARNING Order of execution matters
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
