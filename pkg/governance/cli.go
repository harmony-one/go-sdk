package governance

import (
	"fmt"

	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
)

func DoVote(keyStore *keystore.KeyStore, account accounts.Account, vote Vote) error {
	typedData, err := vote.ToEIP712()
	if err != nil {
		return err
	}
	sig, err := signTypedData(keyStore, account, typedData)
	if err != nil {
		return err
	}

	result, err := submitMessage(account.Address.String(), typedData, sig)
	if err != nil {
		return err
	}

	fmt.Println(indent(result))
	return nil
}
