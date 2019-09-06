package keys

import (
	"fmt"
	"os"
	"path"
	"strings"
	//"crypto"
	//"crypto/ecdsa"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/go-sdk/pkg/common/address"
	"github.com/harmony-one/harmony/accounts/keystore"
	//"github.com/jbenet/go-base58"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	// "github.com/ethereum/go-ethereum/crypto"

	homedir "github.com/mitchellh/go-homedir"
	bip39 "github.com/tyler-smith/go-bip39"
	bip32 "github.com/tyler-smith/go-bip32"
)

func checkAndMakeKeyDirIfNeeded() string {
	userDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(userDir, ".hmy_cli", "keystore")
	if _, err := os.Stat(hmyCLIDir); os.IsNotExist(err) {
		// Double check with Leo what is right file persmission
		os.Mkdir(hmyCLIDir, 0700)
	}

	return hmyCLIDir
}

func ListKeys(keystoreDir string) {
	hmyCLIDir := checkAndMakeKeyDirIfNeeded()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(hmyCLIDir, scryptN, scryptP)
	allAccounts := ks.Accounts()
	fmt.Printf("Harmony Address:%s File URL:\n", strings.Repeat(" ", ethCommon.AddressLength*2))
	for _, account := range allAccounts {
		fmt.Printf("%s\t\t %s\n", address.ToBech32(account.Address), account.URL)
	}

}

func AddNewKey(password string) {
	hmyCLIDir := checkAndMakeKeyDirIfNeeded()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(hmyCLIDir, scryptN, scryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		fmt.Printf("new account error: %v\n", err)
	}
	fmt.Printf("account: %s\n", address.ToBech32(account.Address))
	fmt.Printf("URL: %s\n", account.URL)
}

//Generates a Mnemonic without any input from user like how the WebWallet did it in dapp-examples
func GenerateMnemonic()(string){
	seed, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(seed)
	return mnemonic
}

func NewAccountByMnemonic(mnemonic string)([]byte, []byte) {
	seed, _ := bip39.NewSeedWithErrorChecking(mnemonic, "")
	masterKey, _ := bip32.NewMasterKey(seed)
	/*hmyCLIDir := checkAndMakeKeyDirIfNeeded()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(hmyCLIDir, scryptN, scryptP)*/
	//account, _ := ks.Import(masterKey.Key, "", "")		//can add password functionality here
	//fmt.Printf("account: %s\n", address.ToBech32(account.Address))
	//fmt.Printf("URL: %s\n", account.URL)
	key, _ := masterKey.NewChildKey(2147483648 + 44)
	/*decoded := base58.Decode(key.B58Serialize())
	privateKey := decoded[46:78]*/
	privateKeyECDSA, _ := ethCrypto.ToECDSA(key.Key)
	//privateKeyECDSA, _ := ecdsa.GenerateKey(crypto.S256(), seed)
	return seed, ethCrypto.FromECDSA(privateKeyECDSA)
}