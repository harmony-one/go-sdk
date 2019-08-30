package keys

import (
	"fmt"
	"os"
	"path"
	"strings"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/go-sdk/pkg/common/address"
	"github.com/harmony-one/harmony/accounts/keystore"

	// "github.com/ethereum/go-ethereum/crypto"

	homedir "github.com/mitchellh/go-homedir"
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

func AddNewKey() {
	hmyCLIDir := checkAndMakeKeyDirIfNeeded()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(hmyCLIDir, scryptN, scryptP)
	password := ""
	// TODO Need to factor out some definitions from harmony/internal to something
	// more public, like core or api
	// password := utils.AskForPassphrase("Passphrase: ")
	// password2 := utils.AskForPassphrase("Passphrase again: ")
	// if password != password2 {
	// 	fmt.Printf("Passphrase doesn't match. Please try again!\n")
	// 	os.Exit(3)
	// }
	account, err := ks.NewAccount(password)
	if err != nil {
		fmt.Printf("new account error: %v\n", err)
	}
	// fmt.Printf("URL: %s\n", account.URL)
	fmt.Printf("account: %s\n", address.ToBech32(account.Address))
	fmt.Printf("URL: %s\n", account.URL)
}
