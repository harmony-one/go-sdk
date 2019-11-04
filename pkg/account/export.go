package account

import (
	"fmt"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/harmony/accounts"
)

func ExportPrivateKey(address, passphrase string) error {
	ks := store.FromAddress(address)
	allAccounts := ks.Accounts()
	for _, account := range allAccounts {
		_, key, err := ks.GetDecryptedKey(accounts.Account{Address: account.Address}, passphrase)
		if err != nil {
			return err
		}
		fmt.Printf("%x\n", key.PrivateKey.D)
	}
	return nil
}

func ExportKeystore(address, passphrase string) error {
	ks := store.FromAddress(address)
	allAccounts := ks.Accounts()
	for _, account := range allAccounts {
		keyFile, err := ks.Export(accounts.Account{Address: account.Address}, passphrase, passphrase)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", keyFile)
	}
	return nil
}
