package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/harmony-one/go-sdk/pkg/account"
	c "github.com/harmony-one/go-sdk/pkg/common"

	"github.com/harmony-one/go-sdk/pkg/keys"
	"github.com/harmony-one/go-sdk/pkg/ledger"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	seedPhraseWarning = "**Important** write this seed phrase in a safe place, " +
		"it is the only way to recover your account if you ever forget your password\n\n"
)

var (
	quietImport            bool
	recoverFromMnemonic    bool
	userProvidesPassphrase bool
	importPassphrase       string
	blsFilePath            string
)

func doubleTakePhrase() string {
	fmt.Println("Enter passphrase:")
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
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List all the local accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if useLedgerWallet {
				ledger.ProcessAddressCommand()
				return nil
			}
			store.DescribeLocalAccounts()
			return nil
		},
	}

	cmdLocation := &cobra.Command{
		Use:   "location",
		Short: "Show where `hmy` keeps accounts & their keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(store.DefaultLocation())
			return nil
		},
	}

	cmdAdd := &cobra.Command{
		Use:   "add <ACCOUNT_NAME>",
		Short: "Create a new keystore key",
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
	cmdAdd.Flags().BoolVar(&recoverFromMnemonic, "recover", false, "create keys from a mnemonic")
	ppPrompt := fmt.Sprintf(
		"prompt user for passphrase, otherwise default passphrase: \"`%s`\"", c.DefaultPassphrase,
	)
	cmdAdd.Flags().BoolVar(&userProvidesPassphrase, "use-own-passphrase", false, ppPrompt)

	cmdRemove := &cobra.Command{
		Use:   "remove <ACCOUNT_NAME>",
		Short: "Remove a key from the keystore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !store.DoesNamedAccountExist(args[0]) {
				return fmt.Errorf("account %s doesn't exist", args[0])
			}
			account.RemoveAccount(args[0])
			return nil
		},
	}

	cmdMnemonic := &cobra.Command{
		Use:   "mnemonic",
		Short: "Compute the bip39 mnemonic for some input entropy",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(mnemonic.Generate())
			return nil
		},
	}

	cmdImportKS := &cobra.Command{
		Use:   "import-ks <ABSOLUTE_PATH_KEYSTORE> [ACCOUNT_NAME]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Import an existing keystore key",
		RunE: func(cmd *cobra.Command, args []string) error {
			userName := ""
			if len(args) == 2 {
				userName = args[1]
			}
			name, err := account.ImportKeyStore(args[0], userName, importPassphrase)
			if !quietImport && err == nil {
				fmt.Printf("Imported keystore given account alias of `%s`\n", name)
			}
			return err
		},
	}
	importP := `passphrase of key being imported, default assumes ""`
	cmdImportKS.Flags().StringVar(&importPassphrase, "passphrase", "", importP)
	cmdImportKS.Flags().BoolVar(&quietImport, "quiet", false, "do not print out imported account name")
	cmdImportKS.MarkFlagRequired("passphrase")

	cmdImportSK := &cobra.Command{
		Use:   "import-private-key <secp256k1_PRIVATE_KEY> [ACCOUNT_NAME]",
		Short: "Import an existing keystore key (only accept secp256k1 private keys)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase := c.DefaultPassphrase
			if userProvidesPassphrase {
				passphrase = doubleTakePhrase()
			}
			userName := ""
			if len(args) == 2 {
				userName = args[1]
			}
			name, err := account.ImportFromPrivateKey(args[0], userName, passphrase)
			if !quietImport && err == nil {
				fmt.Printf("Imported keystore given account alias of `%s`\n", name)
			}
			return err
		},
	}
	cmdImportSK.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)

	cmdExportSK := &cobra.Command{
		Use:   "export-private-key <ACCOUNT_ADDRESS>",
		Short: "Export the secp256k1 private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := account.ExportPrivateKey(args[0], unlockP)
			return err
		},
	}
	cmdExportSK.Flags().StringVar(&unlockP,
		"passphrase", c.DefaultPassphrase,
		"passphrase to unlock sender's keystore",
	)

	cmdExportKS := &cobra.Command{
		Use:   "export-ks <ACCOUNT_ADDRESS>",
		Short: "Export the keystore file contents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := account.ExportKeystore(args[0], unlockP)
			return err
		},
	}
	cmdExportKS.Flags().StringVar(&unlockP,
		"passphrase", c.DefaultPassphrase,
		"passphrase to unlock sender's keystore",
	)

	cmdGenerateBlsKey := &cobra.Command{
		Use:   "generate-bls-key",
		Short: "Generate bls keys then encrypt and save the private key with a requested passphrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase := doubleTakePhrase()
			return keys.GenBlsKeys(passphrase, blsFilePath)
		},
	}
	cmdGenerateBlsKey.Flags().StringVar(&blsFilePath, "bls-file-path", "",
		"absolute path of where to save encrypted bls private key")

	cmdRecoverBlsKey := &cobra.Command{
		Use:   "recover-bls-key <ABSOLUTE_PATH_BLS_KEY>",
		Short: "Recover bls keys from an encrypted bls key file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return keys.RecoverBlsKeyFromFile(unlockP, args[0])
		},
	}
	cmdRecoverBlsKey.Flags().StringVar(&unlockP,
		"passphrase", c.DefaultPassphrase,
		"passphrase to unlock sender's keystore",
	)

	cmdSaveBlsKey := &cobra.Command{
		Use:   "save-bls-key <PRIVATE_BLS_KEY>",
		Short: "Encrypt and save the bls private key with a requested passphrase",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase := doubleTakePhrase()
			return keys.SaveBlsKey(passphrase, blsFilePath, args[0])
		},
	}
	cmdSaveBlsKey.Flags().StringVar(&blsFilePath, "bls-file-path", "",
		"absolute path of where to save encrypted bls private key")

	GetPublicBlsKey := &cobra.Command{
		Use:   "get-public-bls-key <PRIVATE_BLS_KEY>",
		Short: "Get the public bls key associated with the provided private bls key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return keys.GetPublicBlsKey(args[0])
		},
	}

	return []*cobra.Command{cmdList, cmdLocation, cmdAdd, cmdRemove, cmdMnemonic, cmdImportKS, cmdImportSK,
		cmdExportKS, cmdExportSK, cmdGenerateBlsKey, cmdRecoverBlsKey, cmdSaveBlsKey, GetPublicBlsKey}
}

func init() {
	cmdKeys := &cobra.Command{
		Use:   "keys",
		Short: "Add or view local private keys",
		Long: `
Manage your local keys
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdKeys.AddCommand(keysSub()...)
	RootCmd.AddCommand(cmdKeys)
}
