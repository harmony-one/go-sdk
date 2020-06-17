package account

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/harmony/accounts"
)

var (
	ErrAddressNotFound = errors.New("account was not found in keystore")
)

func ExportPrivateKey(address, passphrase string) error {
	ks := store.FromAddress(address)
	if ks == nil {
		return ErrAddressNotFound
	}
	account := ks.Accounts()[0]
	_, key, err := ks.GetDecryptedKey(accounts.Account{Address: account.Address}, passphrase)
	if err != nil {
		return err
	}
	fmt.Printf("%064x\n", key.PrivateKey.D)
	return nil
}

func VerifyPassphrase(address, passphrase string) (bool, error) {
	ks := store.FromAddress(address)
	if ks == nil {
		return false, ErrAddressNotFound
	}
	account := ks.Accounts()[0]
	_, _, err := ks.GetDecryptedKey(accounts.Account{Address: account.Address}, passphrase)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExportKeystore(address, path, passphrase string) (string, error) {
	ks := store.FromAddress(address)
	if ks == nil {
		return "", ErrAddressNotFound
	}
	account := ks.Accounts()[0]
	dirPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	outFile := filepath.Join(dirPath, fmt.Sprintf("%s.key", address))
	keyFile, err := ks.Export(accounts.Account{Address: account.Address}, passphrase, passphrase)
	if err != nil {
		return "", err
	}
	e := writeToFile(outFile, string(keyFile))
	if e != nil {
		return "", e
	}
	return outFile, nil
}
