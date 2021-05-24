package console

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/console/jsre"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/crypto/hash"
	"math/big"
	"strconv"
	"time"
)

func getStringFromJsObjWithDefault(o *goja.Object, key string, def string) string {
	get := o.Get(key)
	if get == nil {
		return def
	} else {
		return get.String()
	}
}

func (b *bridge) callbackProtected(protectedFunc func(call jsre.Call) (goja.Value, error)) func(call jsre.Call) (goja.Value, error) {
	return func(call jsre.Call) (goja.Value, error) {
		var availableCB goja.Callable = nil
		for i, args := range call.Arguments {
			if cb, ok := goja.AssertFunction(args); ok {
				availableCB = cb
				call.Arguments = call.Arguments[:i] // callback must be last
				break
			}
		}

		value, err := protectedFunc(call)
		jsErr := goja.Undefined()
		if err != nil {
			jsErr = call.VM.NewGoError(err)
		}
		if availableCB != nil {
			_, _ = availableCB(nil, jsErr, value)
		}

		return value, err
	}
}

func (b *bridge) HmyGetListAccounts(call jsre.Call) (goja.Value, error) {
	var accounts = []string{}

	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			accounts = append(accounts, account.Address.String())
		}
	}

	return call.VM.ToValue(accounts), nil
}

func (b *bridge) HmySignTransaction(call jsre.Call) (goja.Value, error) {
	txObj := call.Arguments[0].ToObject(call.VM)
	password := call.Arguments[1].String()

	from := getStringFromJsObjWithDefault(txObj, "from", "")
	to := getStringFromJsObjWithDefault(txObj, "to", "")
	gasLimit := getStringFromJsObjWithDefault(txObj, "gas", "1000000")
	amount := getStringFromJsObjWithDefault(txObj, "value", "0")
	gasPrice := getStringFromJsObjWithDefault(txObj, "gasPrice", "1")

	networkHandler := rpc.NewHTTPHandler(b.console.nodeUrl)
	chanId, err := common.StringToChainID(b.console.net)
	if err != nil {
		return nil, err
	}

	ks, acct, err := store.UnlockedKeystore(from, password)
	if err != nil {
		return nil, err
	}
	ctrlr := transaction.NewController(networkHandler, ks, acct, *chanId, func(controller *transaction.Controller) {
		// nop
	})

	tempLimit, err := strconv.ParseInt(gasLimit, 10, 64)
	if err != nil {
		return nil, err
	}
	if tempLimit < 0 {
		return nil, errors.New(fmt.Sprintf("gas-limit can not be negative: %s", gasLimit))
	}
	gLimit := uint64(tempLimit)

	amt, err := common.NewDecFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("amount %w", err)
	}

	gPrice, err := common.NewDecFromString(gasPrice)
	if err != nil {
		return nil, fmt.Errorf("gas-price %w", err)
	}

	toP := &to
	if to == "" {
		toP = nil
	}

	nonce := transaction.GetNextPendingNonce(from, networkHandler)
	err = ctrlr.SignTransaction(
		nonce, gLimit,
		toP,
		uint32(b.console.shardId), uint32(b.console.shardId),
		amt, gPrice,
		[]byte{},
	)
	if err != nil {
		return nil, err
	}

	info := ctrlr.TransactionInfo()

	return call.VM.ToValue(map[string]interface{}{
		"raw": ctrlr.RawTransaction(),
		"tx": map[string]string{
			"nonce":    "0x" + big.NewInt(int64(info.Nonce())).Text(16),
			"gasPrice": "0x" + info.GasPrice().Text(16),
			"gas":      "0x" + big.NewInt(int64(info.GasLimit())).Text(16),
			"to":       info.To().Hex(),
			"value":    "0x" + info.Value().Text(16),
			"input":    "0x" + hex.EncodeToString(info.Data()),
			"v":        "0x" + info.V().Text(16),
			"r":        "0x" + info.R().Text(16),
			"s":        "0x" + info.S().Text(16),
			"hash":     info.Hash().Hex(),
		},
	}), nil
}

func (b *bridge) HmySendTransaction(call jsre.Call) (goja.Value, error) {
	txObj := call.Arguments[0].ToObject(call.VM)
	password := ""
	if len(call.Arguments) > 1 {
		password = call.Arguments[1].String()
	}

	from := getStringFromJsObjWithDefault(txObj, "from", "")
	to := getStringFromJsObjWithDefault(txObj, "to", "")
	gasLimit := getStringFromJsObjWithDefault(txObj, "gas", "1000000")
	amount := getStringFromJsObjWithDefault(txObj, "value", "0")
	gasPrice := getStringFromJsObjWithDefault(txObj, "gasPrice", "1")

	networkHandler := rpc.NewHTTPHandler(b.console.nodeUrl)
	chanId, err := common.StringToChainID(b.console.net)
	if err != nil {
		return nil, err
	}

	ks, acct, err := store.UnlockedKeystore(from, password)
	if err != nil {
		return nil, err
	}
	ctrlr := transaction.NewController(networkHandler, ks, acct, *chanId, func(controller *transaction.Controller) {
		// nop
	})

	tempLimit, err := strconv.ParseInt(gasLimit, 10, 64)
	if err != nil {
		return nil, err
	}
	if tempLimit < 0 {
		return nil, errors.New(fmt.Sprintf("gas-limit can not be negative: %s", gasLimit))
	}
	gLimit := uint64(tempLimit)

	amt, err := common.NewDecFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("amount %w", err)
	}

	gPrice, err := common.NewDecFromString(gasPrice)
	if err != nil {
		return nil, fmt.Errorf("gas-price %w", err)
	}

	toP := &to
	if to == "" {
		toP = nil
	}

	nonce := transaction.GetNextPendingNonce(from, networkHandler)
	err = ctrlr.ExecuteTransaction(
		nonce, gLimit,
		toP,
		uint32(b.console.shardId), uint32(b.console.shardId),
		amt, gPrice,
		[]byte{},
	)
	if err != nil {
		return nil, err
	}

	return call.VM.ToValue(*ctrlr.TransactionHash()), nil
}

func (b *bridge) HmyLockAccount(call jsre.Call) (goja.Value, error) {
	address := call.Arguments[0].String()

	_, _, err := store.LockKeystore(address)
	if err != nil {
		return nil, err
	}

	return goja.Null(), nil
}

func (b *bridge) HmyImportRawKey(call jsre.Call) (goja.Value, error) {
	privateKey := call.Arguments[0].String()
	password := call.Arguments[1].String()

	name, err := account.ImportFromPrivateKey(privateKey, "", password)
	if err != nil {
		return nil, err
	}

	return call.VM.ToValue(name), nil
}

func (b *bridge) HmyUnlockAccount(call jsre.Call) (goja.Value, error) {
	if len(call.Arguments) < 3 {
		return nil, errors.New("arguments < 3")
	}
	address := call.Arguments[0].String()
	password := call.Arguments[1].String()
	unlockDuration := call.Arguments[2].ToInteger()

	_, _, err := store.UnlockedKeystoreTimeLimit(address, password, time.Duration(unlockDuration)*time.Second)
	if err != nil {
		return nil, err
	}

	return goja.Null(), nil
}

func (b *bridge) HmyNewAccount(call jsre.Call) (goja.Value, error) {
	return goja.Null(), nil
}

func (b *bridge) HmySign(call jsre.Call) (goja.Value, error) {
	dataToSign := call.Arguments[0].String()
	addressStr := call.Arguments[1].String()
	password := call.Arguments[2].String()

	ks := store.FromAddress(addressStr)
	if ks == nil {
		return nil, fmt.Errorf("could not open local keystore for %s", addressStr)
	}

	acc, err := ks.Find(accounts.Account{Address: address.Parse(addressStr)})
	if err != nil {
		return nil, err
	}

	message, err := signMessageWithPassword(ks, acc, password, []byte(dataToSign))
	if err != nil {
		return nil, err
	}

	return call.VM.ToValue(hex.EncodeToString(message)), nil
}

func signMessageWithPassword(keyStore *keystore.KeyStore, account accounts.Account, password string, data []byte) (sign []byte, err error) {
	signData := append([]byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(data))))
	msgHash := hash.Keccak256(append(signData, data...))

	sign, err = keyStore.SignHashWithPassphrase(account, password, msgHash)
	if err != nil {
		return nil, err
	}

	if len(sign) != 65 {
		return nil, fmt.Errorf("sign error")
	}

	sign[64] += 0x1b
	return sign, nil
}
