package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	c "github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/pkg/errors"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	uDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	if _, err := os.Stat(hmyCLIDir); os.IsNotExist(err) {
		os.MkdirAll(hmyCLIDir, 0700)
	}
}

// LocalAccounts returns a slice of local account alias names
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

var (
	describe              = fmt.Sprintf("%-24s\t\t%23s\n", "NAME", "ADDRESS")
	NoUnlockBadPassphrase = errors.New("could not unlock wallet with given passphrase")
)

// DescribeLocalAccounts will display all the account alias name and their corresponding one address
func DescribeLocalAccounts() {
	fmt.Println(describe)
	for _, name := range LocalAccounts() {
		ks := FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			fmt.Printf("%-48s\t%s\n", name, address.ToBech32(account.Address))
		}
	}
}

// DoesNamedAccountExist return true if the given string name is an alias account already define,
// and return false otherwise
func DoesNamedAccountExist(name string) bool {
	for _, account := range LocalAccounts() {
		if account == name {
			return true
		}
	}
	return false
}

// Returns one address for account name if exists
func AddressFromAccountName(name string) (string, error) {
	ks := FromAccountName(name)
	// FIXME: Assume 1 account per keystore for now
	for _, account := range ks.Accounts() {
		return address.ToBech32(account.Address), nil
	}
	return "", fmt.Errorf("Keystore not found.")
}

// FromAddress will return nil if the bech32 string is not found in the imported accounts
func FromAddress(bech32 string) *keystore.KeyStore {
	for _, name := range LocalAccounts() {
		ks := FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			if bech32 == address.ToBech32(account.Address) {
				return ks
			}
		}
	}
	return nil
}

func FromAccountName(name string) *keystore.KeyStore {
	uDir, _ := homedir.Dir()
	p := path.Join(uDir, c.DefaultConfigDirName, c.DefaultConfigAccountAliasesDirName, name)
	return common.KeyStoreForPath(p)
}

func DefaultLocation() string {
	uDir, _ := homedir.Dir()
	return path.Join(uDir, c.DefaultConfigDirName, c.DefaultConfigAccountAliasesDirName)
}

func UnlockedKeystore(from, passphrase string) (*keystore.KeyStore, *accounts.Account, error) {
	return UnlockedKeystoreTimeLimit(from, passphrase, 0)
}

func LockKeystore(from string) (*keystore.KeyStore, *accounts.Account, error) {
	sender := address.Parse(from)
	ks := FromAddress(address.ToBech32(sender))
	if ks == nil {
		return nil, nil, fmt.Errorf("could not open local keystore for %s", from)
	}
	account, lookupErr := ks.Find(accounts.Account{Address: sender})
	if lookupErr != nil {
		return nil, nil, fmt.Errorf("could not find %s in keystore", from)
	}
	if lockError := ks.Lock(account.Address); lockError != nil {
		return nil, nil, lockError
	}
	return ks, &account, nil
}

func UnlockedKeystoreTimeLimit(from, passphrase string, time time.Duration) (*keystore.KeyStore, *accounts.Account, error) {
	sender := address.Parse(from)
	ks := FromAddress(address.ToBech32(sender))
	if ks == nil {
		return nil, nil, fmt.Errorf("could not open local keystore for %s", from)
	}
	account, lookupErr := ks.Find(accounts.Account{Address: sender})
	if lookupErr != nil {
		return nil, nil, fmt.Errorf("could not find %s in keystore", from)
	}
	if unlockError := ks.TimedUnlock(account, passphrase, time); unlockError != nil {
		return nil, nil, errors.Wrap(NoUnlockBadPassphrase, unlockError.Error())
	}
	return ks, &account, nil
}
