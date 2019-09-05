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
