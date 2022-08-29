package governance

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/crypto/hash"
	"github.com/pkg/errors"
)

func encodeForSigning(typedData *TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"cannot hash the domain structure",
		)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"cannot hash the structure",
		)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return rawData, nil
}

// signTypedData encodes and signs EIP-712 data
// it is copied over here from Geth to use our own keystore implementation
func signTypedData(keyStore *keystore.KeyStore, account accounts.Account, typedData *TypedData) (string, error) {
	rawData, err := encodeForSigning(typedData)
	if err != nil {
		return "", errors.Wrapf(
			err,
			"cannot encode for signing",
		)
	}

	msgHash := hash.Keccak256Hash(rawData)
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

// func signMessage(keyStore *keystore.KeyStore, account accounts.Account, data []byte) (string, error) {
// 	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
// 	msgHash := hash.Keccak256Hash([]byte(fullMessage))
// 	sign, err := keyStore.SignHash(account, msgHash.Bytes())
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(sign) != 65 {
// 		return "", fmt.Errorf("sign error")
// 	}
// 	sign[64] += 0x1b
// 	return hexutil.Encode(sign), nil
// }
