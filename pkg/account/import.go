package account

import (
	"encoding/hex"
	"io/ioutil"
	"path"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	mapset "github.com/deckarep/golang-set"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/pkg/errors"
)

var (
	// ErrNotAbsPath when keypath not absolute path
	ErrNotAbsPath   = errors.New("keypath is not absolute path")
	ErrBadKeyLength = errors.New("Invalid private key (wrong length)")
)

// ImportFromPrivateKey allows import of an ECDSA private key
func ImportFromPrivateKey(privateKey, name, passphrase string) (string, error) {
	if name == "" {
		name = generateName() + "-imported"
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	if len(privateKeyBytes) != common.Secp256k1PrivateKeyBytesLength {
		return "", ErrBadKeyLength
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
func ImportKeyStore(keypath, name, passphrase string) (string, error) {
	if !path.IsAbs(keypath) {
		return "", ErrNotAbsPath
	}
	keyJSON, readError := ioutil.ReadFile(keypath)
	if readError != nil {
		return "", readError
	}
	if name == "" {
		name = generateName() + "-imported"
	}
	ks := store.FromAccountName(name)
	_, err := ks.Import(keyJSON, passphrase, common.DefaultPassphrase)
	if err != nil {
		return "", errors.Wrap(err, "could not import")
	}

	return name, nil
}
