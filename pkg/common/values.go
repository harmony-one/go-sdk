package common

import (
	"errors"
	"os"

	"github.com/harmony-one/harmony/accounts/keystore"
)

const (
	DefaultConfigDirName               = ".hmy_cli"
	DefaultConfigAccountAliasesDirName = "account-keys"
	DefaultPassphrase                  = ""
	JSONRPCVersion                     = "2.0"
	Secp256k1PrivateKeyBytesLength     = 32
)

var (
	ScryptN          = keystore.StandardScryptN
	ScryptP          = keystore.StandardScryptP
	DebugRPC         = false
	DebugTransaction = false
	ErrNotAbsPath    = errors.New("keypath is not absolute path")
	ErrBadKeyLength  = errors.New("Invalid private key (wrong length)")
)

func init() {
	if _, enabled := os.LookupEnv("HMY_RPC_DEBUG"); enabled != false {
		DebugRPC = true
	}
	if _, enabled := os.LookupEnv("HMY_TX_DEBUG"); enabled != false {
		DebugTransaction = true
	}
	if _, enabled := os.LookupEnv("HMY_ALL_DEBUG"); enabled != false {
		EnableAllVerbose()
	}
}

// EnableAllVerbose sets debug vars to true
func EnableAllVerbose() {
	DebugRPC = true
	DebugTransaction = true
}
