package account

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/pkg/errors"
)

var (
	// ErrNotAbsPath when keypath not absolute path
	ErrNotAbsPath = errors.New("keypath is not absolute path")
)

// ImportFromPrivateKey allows import of an ECDSA private key
func ImportFromPrivateKey(privateKey []byte) (string, error) {
	name := generateName() + "-imported"
	ks := store.FromAccountName(name)
	sk, b := x509.ParseECPrivateKey(privateKey)
	// secp256k1
	/*
				openssl ecparam -genkey -name secp256k1 -text -noout -outform DER | xxd -p -c 1000 | sed 's/41534e31204f49443a20736563703235366b310a30740201010420/PrivKey: /' | sed 's/a00706052b8104000aa144034200/\'$'\nPubKey: /'


				openssl ecparam -genkey -name secp256k1 -text -noout -outform DER | xxd -p -c 1000 | sed 's/41534e31204f49443a20736563703235366b310a30740201010420//' | sed 's/a00706052b8104000aa144034200/\'$'\n/'

				 openssl ecparam -name secp256k1  -param_enc explicit -genkey -out key.pem


				https://developers.yubico.com/PIV/Guides/Generating_keys_using_OpenSSL.html

				https://bitcoin.stackexchange.com/questions/59644/how-do-these-openssl-commands-create-a-bitcoin-private-key-from-a-ecdsa-keypair

				http://www.herongyang.com/EC-Cryptography/EC-Key-secp256k1-with-OpenSSL.html

				https://crypto.stackexchange.com/questions/50019/public-key-format-for-ecdsa-as-in-fips-186-4

				https://security.stackexchange.com/questions/84327/converting-ecc-private-key-to-pkcs1-format

				https://serverfault.com/questions/9708/what-is-a-pem-file-and-how-does-it-differ-from-other-openssl-generated-key-file

				https://golang.org/pkg/crypto/x509/#ParseECPrivateKey

				https://github.com/ethereum/go-ethereum/blob/master/vendor/golang.org/x/crypto/ssh/keys.go

				https://godoc.org/golang.org/x/crypto/ssh

		openssl ecparam -list_curves

	*/

	a, err := ks.ImportECDSA(sk, common.DefaultPassphrase)

	fmt.Println(sk, b, a, err)

	if err != nil {
		return "", err
	}
	return "", nil

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
func ImportKeyStore(keypath, passphrase string) (string, error) {
	if !path.IsAbs(keypath) {
		return "", ErrNotAbsPath
	}
	keyJSON, readError := ioutil.ReadFile(keypath)
	if readError != nil {
		return "", readError
	}
	name := generateName() + "-imported"
	ks := store.FromAccountName(name)
	_, err := ks.Import(keyJSON, passphrase, common.DefaultPassphrase)
	if err != nil {
		return "", errors.Wrap(err, "could not import")
	}

	return name, nil
}
