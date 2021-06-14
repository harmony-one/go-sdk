package console

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/crypto/hash"
	"strconv"
)

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

func getStringFromJsObjWithDefault(o *goja.Object, key string, def string) string {
	get := o.Get(key)
	if get == nil {
		return def
	} else {
		return get.String()
	}
}
