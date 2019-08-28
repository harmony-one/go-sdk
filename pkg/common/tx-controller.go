package common

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"os"

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
)

type TxController struct {
	failure          error
	messenger        rpc.T
	node             string
	useOneAddressRPC bool
	ks               *keystore.KeyStore
}

type rpcReply map[string]interface{}

// TODO Make node more robust with URL validation
func NewTxController(handler rpc.T, node string, useOneAddressInsteadOfHex bool) *TxController {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	return &TxController{
		failure:          nil,
		messenger:        handler,
		node:             node,
		useOneAddressRPC: useOneAddressInsteadOfHex, // TODO Not hard coded but parameterized
		ks:               keystore.NewKeyStore("/Users/edgar/.hmy_cli/keystore", scryptN, scryptP),
	}
}

func (Controller *TxController) balance(params []interface{}) rpcReply {
	return rpc.RPCRequest(rpc.Method.GetBalance, Controller.node, params)
}

func (Controller *TxController) transactionCount(params []interface{}) rpcReply {
	return rpc.RPCRequest(rpc.Method.GetTransactionCount, Controller.node, params)
}

func (Controller *TxController) sendSignedRawTx(params []interface{}) rpcReply {
	return rpc.RPCRequest(rpc.Method.SendRawTransaction, Controller.node, params)
}

func (Controller *TxController) txReceipt(params []interface{}) rpcReply {
	return rpc.RPCRequest(rpc.Method.GetTransactionReceipt, Controller.node, params)
}

// func DoTransaction

// TODO Respect the .useOneAddressRPC field for when actually sending it off
// Get current transaction count, that's your new nonce
// Get balance
// Get gas issue
// Then kick it off

func (Controller *TxController) CreateTransaction(
	from, to string,
	amount float64,
	fromShard, toShard int) []byte {
	senderAddress := address.Parse(from)
	receiverAddress := address.Parse(to)
	transactionCountRPCReply := Controller.transactionCount([]interface{}{senderAddress.Hex(), "latest"})

	// TODO Handle the failure case or be more sure about result being fine
	transactionCount, _ := transactionCountRPCReply["result"].(string)

	nonce, _ := big.NewInt(0).SetString(transactionCount[2:], 16)

	// TODO Why the latest param? I forgot, used to know
	balanceRPCReply := Controller.balance([]interface{}{address.ToBech32(senderAddress), "latest"})
	currentBalance, _ := balanceRPCReply["result"].(string)
	balance := big.NewInt(0)
	// TODO Not sure if better to index like this or use the Replace function
	// n, _ := balance.SetString(strings.Replace(currentBalance, "0x", "", -1), 16)
	balance, _ = balance.SetString(currentBalance[2:], 16)
	// fmt.Println(ConvertBalanceIntoReadableFormat(balance))

	account, err := Controller.ks.Find(accounts.Account{Address: senderAddress})

	// fmt.Println(account, err)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return nil
	}
	// TODO Smart way to unlock account, think with John
	Controller.ks.Unlock(account, "edgar")
	amountBigInt := big.NewInt(int64(amount * denominations.Nano))
	amountBigInt = amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))
	inputData, _ := base64.StdEncoding.DecodeString("")
	gas, _ := core.IntrinsicGas(inputData, false, true)
	// TODO Refactor to use the cross-shard transaction item
	// tx := hmyTypes.NewCrossShardTransaction(
	// 	transactionCount+1, &receiverAddress, fromShard, toShard, amountBigInt,
	// 	gas, gasPriceBigInt, inputData)

	// fmt.Println(nonce.Uint64())
	tx := types.NewTransaction(
		nonce.Uint64(), receiverAddress, uint32(0), amountBigInt,
		gas, nil, inputData)

	r, _ := tx.MarshalJSON()
	fmt.Println(string(r))

	signedTransaction, _ := Controller.ks.SignTx(account, tx, nil)

	// ts := types.Transactions{signedTransaction}
	// rawTx := hexutil.Encode(ts.GetRlp(0))

	enc, _ := rlp.EncodeToBytes(signedTransaction)
	rawTx := hexutil.Encode(enc)

	sendRawTransactionRPCReply := Controller.sendSignedRawTx([]interface{}{rawTx})
	txReceipt, _ := sendRawTransactionRPCReply["result"].(string)

	fmt.Println(sendRawTransactionRPCReply, txReceipt)

	txReceiptRPCReply := Controller.txReceipt([]interface{}{txReceipt})
	fmt.Println(txReceiptRPCReply)
	return nil
}

func (Controller *TxController) SignTransaction(arg []byte) []byte {
	return nil
}

func (Controller *TxController) SendTransaction(arg []byte) []byte {
	return nil
}
