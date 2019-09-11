package common

import (
	"github.com/harmony-one/harmony/accounts/keystore"
)

const (
	DefaultConfigDirName               = ".hmy_cli"
	DefaultConfigAccountAliasesDirName = "account-keys"
	DefaultPassphrase                  = "harmony-one"
)

var (
	ScryptN = keystore.StandardScryptN
	ScryptP = keystore.StandardScryptP
)
