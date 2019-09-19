package cmd

import (
	"bufio"
	"fmt"
	"os"

	color "github.com/fatih/color"
	"github.com/harmony-one/go-sdk/pkg/account"
	c "github.com/harmony-one/go-sdk/pkg/common"

	"github.com/harmony-one/go-sdk/pkg/ledger"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	seedPhraseWarning = ("**Important** write this seed phrase in a safe place, " +
		"it is the only way to recover your account if you ever forget your password\n\n")
)

var (
	quietImport            bool
	recoverFromMnemonic    bool
	userProvidesPassphrase bool
	importPassphrase       string
)

func doubleTakePhrase() string {
	fmt.Println("Enter passphrase for account")
	pass, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("Repeat the passphrase:")
	repeatPass, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	if string(repeatPass) != string(pass) {
		fmt.Println("Passphrases do not match")
		os.Exit(-1)
	}
	return string(repeatPass)
}

func keysSub() []*cobra.Command {
	add := &cobra.Command{
		Use:   "add <ACCOUNT_NAME>",
		Short: "Create a new key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if store.DoesNamedAccountExist(args[0]) {
				return fmt.Errorf("account %s already exists", args[0])
			}
			passphrase := c.DefaultPassphrase
			if userProvidesPassphrase {
				passphrase = doubleTakePhrase()
			}
			t := account.Creation{args[0], passphrase, "", nil, nil}
			if recoverFromMnemonic {
				fmt.Println("Enter mnemonic to recover keys from")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				m := scanner.Text()
				if !bip39.IsMnemonicValid(m) {
					return mnemonic.InvalidMnemonic
				}
				t.Mnemonic = m
			}
			if err := account.CreateNewLocalAccount(&t); err != nil {
				return err
			}
			if !recoverFromMnemonic {
				color.Red(seedPhraseWarning)
				fmt.Println(t.Mnemonic)
			}
			return nil
		},
	}
	// add.Flags().BoolVar(&recoverFromMnemonic, "recover", false, "create keys from a mnemonic")
	ppPrompt := fmt.Sprintf("provide own phrase over default: `%s`", c.DefaultPassphrase)
	add.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdImport := &cobra.Command{
		Use:   "import <ABSOLUTE_PATH_KEYSTORE>",
		Args:  cobra.ExactArgs(1),
		Short: "Import an existing keystore key",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := account.ImportKeyStore(args[0], importPassphrase)
			if !quietImport {
				fmt.Printf("Imported keystore given account alias of `%s`\n", name)
			}
			return err
		},
	}
	importP := `passphrase of key being imported, default assumes ""`
	cmdImport.Flags().StringVar(&importPassphrase, "passphrase", "", importP)
	cmdImport.Flags().BoolVar(&quietImport, "quiet", false, "do not print out imported account name")

	return []*cobra.Command{add, cmdImport, {
		Use:   "mnemonic",
		Short: "Compute the bip39 mnemonic for some input entropy",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(mnemonic.Generate())
		},
	}, {
		Use:   "list",
		Short: "List all the local accounts",
		Run: func(cmd *cobra.Command, args []string) {
			if useLedgerWallet {
				ledger.ProcessAddressCommand()
				return
			}
			store.DescribeLocalAccounts()
		},
	}, {
		Use:   "location",
		Short: "Show where `hmy` keeps accounts & their keys",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(store.DefaultLocation())
		},
	},
	}
}

func init() {
	cmdKeys := &cobra.Command{
		Use:   "keys",
		Short: "Add or view local private keys",
		Long: `
Manage your local keys
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmdKeys.AddCommand(keysSub()...)
	RootCmd.AddCommand(cmdKeys)
}
