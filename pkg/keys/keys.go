package keys

import (
	"fmt"
	"os"
	"path"

	"github.com/harmony-one/harmony/accounts/keystore"
	// common2 "github.com/harmony-one/harmony/internal/common"
	// "github.com/harmony-one/harmony/internal/utils"

	homedir "github.com/mitchellh/go-homedir"
)

func checkAndMakeKeyDirIfNeeded() string {
	userDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(userDir, ".hmy_cli")
	if _, err := os.Stat(hmyCLIDir); os.IsNotExist(err) {
		// Double check with Leo what is right file persmission
		os.Mkdir(hmyCLIDir, 0700)
	}

	return hmyCLIDir
}

func ListKeys() {
	hmyCLIDir := checkAndMakeKeyDirIfNeeded()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(hmyCLIDir, scryptN, scryptP)
	allAccounts := ks.Accounts()
	for _, account := range allAccounts {
		fmt.Printf("account: %s\n", account.Address)
		fmt.Printf("URL: %s\n", account.URL)
	}

}

func AddNewKey() {
	hmyCLIDir := checkAndMakeKeyDirIfNeeded()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(hmyCLIDir, scryptN, scryptP)
	password := "edgar"
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
	// fmt.Printf("account: %s\n", common2.MustAddressToBech32(account.Address))
	// fmt.Printf("URL: %s\n", account.URL)
	fmt.Printf("account: %s\n", account.Address)
	fmt.Printf("URL: %s\n", account.URL)
}
