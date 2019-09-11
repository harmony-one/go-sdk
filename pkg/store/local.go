package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	c "github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/accounts/keystore"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	uDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	if _, err := os.Stat(hmyCLIDir); os.IsNotExist(err) {
		os.MkdirAll(hmyCLIDir, 0700)
	}
}

func LocalAccounts() []string {
	uDir, _ := homedir.Dir()
	files, _ := ioutil.ReadDir(path.Join(
		uDir,
		common.DefaultConfigDirName,
		common.DefaultConfigAccountAliasesDirName,
	))
	accounts := []string{}
	for _, node := range files {
		if node.IsDir() {
			accounts = append(accounts, path.Base(node.Name()))
		}
	}
	return accounts
}

func DescribeLocalAccounts() {
	fmt.Printf("NAME\t\tADDRESS\n")
	for _, name := range LocalAccounts() {
		ks := FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			fmt.Printf("%s\t\t %s\n", name, address.ToBech32(account.Address))
		}
	}
}

func DoesNamedAccountExist(name string) bool {
	for _, account := range LocalAccounts() {
		if account == name {
			return true
		}
	}
	return false
}

func FromAccountName(name string) *keystore.KeyStore {
	uDir, _ := homedir.Dir()
	p := path.Join(uDir, c.DefaultConfigDirName, c.DefaultConfigAccountAliasesDirName, name)
	return common.KeyStoreForPath(p)
}
