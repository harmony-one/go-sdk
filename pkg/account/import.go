package account

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"os"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/harmony/accounts/keystore"
	c "github.com/harmony-one/go-sdk/pkg/common"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/btcsuite/btcd/btcec"
	mapset "github.com/deckarep/golang-set"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	"github.com/harmony-one/go-sdk/pkg/store"
)

// ImportFromPrivateKey allows import of an ECDSA private key
func ImportFromPrivateKey(privateKey, name, passphrase string) (string, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")

	if name == "" {
		name = generateName() + "-imported"
		for store.DoesNamedAccountExist(name) {
			name = generateName() + "-imported"
		}
	} else if store.DoesNamedAccountExist(name) {
		return "", fmt.Errorf("account %s already exists", name)
	}

	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	if len(privateKeyBytes) != common.Secp256k1PrivateKeyBytesLength {
		return "", common.ErrBadKeyLength
	}

	// btcec.PrivKeyFromBytes only returns a secret key and public key
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	ks := store.FromAccountName(name)
	_, err = ks.ImportECDSA(sk.ToECDSA(), passphrase)
	return name, err
}

func generateName() string {
	words := strings.Split(mnemonic.Generate(), " ")
	existingAccounts := mapset.NewSet()
	for a := range store.LocalAccounts() {
		existingAccounts.Add(a)
	}
	foundName := false
	acct := ""
	i := 0
	for {
		if foundName {
			break
		}
		if i == len(words)-1 {
			words = strings.Split(mnemonic.Generate(), " ")
		}
		candidate := words[i]
		if !existingAccounts.Contains(candidate) {
			foundName = true
			acct = candidate
			break
		}
	}
	return acct
}

// ImportKeyStore imports a keystore along with a password
func ImportKeyStore(keyPath, name, passphrase string) (string, error) {
	keyPath, err := filepath.Abs(keyPath)
	if err != nil {
		return "", err
	}
	keyJSON, readError := ioutil.ReadFile(keyPath)
	if readError != nil {
		return "", readError
	}
	if name == "" {
		name = generateName() + "-imported"
		for store.DoesNamedAccountExist(name) {
			name = generateName() + "-imported"
		}
	} else if store.DoesNamedAccountExist(name) {
		return "", fmt.Errorf("account %s already exists", name)
	}
	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return "", err
	}
	b32 := address.ToBech32(key.Address)
	hasAddress := store.FromAddress(b32) != nil
	if hasAddress {
		return "", fmt.Errorf("address %s already exists in keystore", b32)
	}
	uDir, _ := homedir.Dir()
	path := filepath.Join(uDir, c.DefaultConfigDirName, c.DefaultConfigAccountAliasesDirName, name)	
	err = os.MkdirAll(path, 0600)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(path, filepath.Base(keyPath)), keyJSON, 0600)
	if err != nil {
		return "", err
	}
	return name, nil
}
