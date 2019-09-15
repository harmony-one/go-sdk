package account

import (
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
	NotAbsPath = errors.New("keypath is not absolute path")
)

func ImportKeyStore(keypath string) error {
	if !path.IsAbs(keypath) {
		return NotAbsPath
	}
	keyJSON, readError := ioutil.ReadFile(keypath)
	if readError != nil {
		return readError
	}
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
	ks := store.FromAccountName(acct + "-imported")
	_, err := ks.Import(keyJSON, "", common.DefaultPassphrase)
	if err != nil {
		return errors.Wrap(err, "could not import")
	}

	return nil
}
