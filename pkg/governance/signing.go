package governance

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/crypto/hash"
)

func signMessage(keyStore *keystore.KeyStore, account accounts.Account, data []byte) (string, error) {
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	msgHash := hash.Keccak256Hash([]byte(fullMessage))
	sign, err := keyStore.SignHash(account, msgHash.Bytes())
	if err != nil {
		return "", err
	}
	if len(sign) != 65 {
		return "", fmt.Errorf("sign error")
	}
	sign[64] += 0x1b
	return hexutil.Encode(sign), nil
}
