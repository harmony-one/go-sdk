package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
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
	userProvidesPassphrase bool
	passphraseFilePath     string
	passphrase             string
	blsFilePath            string
	ppPrompt               = fmt.Sprintf(
		"prompt for passphrase, otherwise use default passphrase: \"`%s`\"", c.DefaultPassphrase,
	)
)

// getPassphrase fetches the correct passphrase depending on if a file is available to
// read from or if the user wants to enter in their own passphrase. Otherwise, just use
// the default passphrase.
func getPassphrase() (string, error) {
	if passphraseFilePath != "" {
		if _, err := os.Stat(passphraseFilePath); os.IsNotExist(err) {
			return "", errors.New(fmt.Sprintf("passphrase file not found at `%s`", passphraseFilePath))
		}
		dat, err := ioutil.ReadFile(passphraseFilePath)
		if err != nil {
			return "", err
		}
		return string(dat), nil
	} else if userProvidesPassphrase {
		fmt.Println("Enter passphrase:")
		pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		fmt.Println("Repeat the passphrase:")
		repeatPass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		if string(repeatPass) != string(pass) {
			return "", errors.New("passphrase does not match")
		}
		fmt.Println("") // provide feedback when passphrase is entered.
		return string(repeatPass), nil
	} else {
		return c.DefaultPassphrase, nil
	}
}

func keysSub() []*cobra.Command {
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List all the local accounts",
		RunE:  func(cmd *cobra.Command, args []string) error {
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
		RunE:  func(cmd *cobra.Command, args []string) error {
			fmt.Println(store.DefaultLocation())
			return nil
		},
	}

	cmdAdd := &cobra.Command{
		Use:   "add <ACCOUNT_NAME>",
		Short: "Create a new keystore key",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error {
			if store.DoesNamedAccountExist(args[0]) {
				return fmt.Errorf("account %s already exists", args[0])
			}
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			acc := account.Creation{
				Name:       args[0],
				Passphrase: passphrase,
			}
			if err := account.CreateNewLocalAccount(&acc); err != nil {
				return err
			}
			color.Red(seedPhraseWarning)
			fmt.Println(acc.Mnemonic)
			return nil
		},
	}
	cmdAdd.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdAdd.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	cmdRemove := &cobra.Command{
		Use:   "remove <ACCOUNT_NAME>",
		Short: "Remove a key from the keystore",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error {
			if err := account.RemoveAccount(args[0]); err != nil {
				return err
			}
			return nil
		},
	}

	cmdMnemonic := &cobra.Command{
		Use:   "mnemonic",
		Short: "Compute the bip39 mnemonic for some input entropy",
		RunE:  func(cmd *cobra.Command, args []string) error {
			fmt.Println(mnemonic.Generate())
			return nil
		},
	}

	cmdRecoverMnemonic := &cobra.Command{
		Use:   "recover-from-mnemonic [ACCOUNT_NAME]",
		Short: "Recover account from mnemonic",
		Args:  cobra.ExactArgs(1),
		RunE:  func (cmd *cobra.Command, args []string) error {
			if store.DoesNamedAccountExist(args[0]) {
	      return fmt.Errorf("account %s already exists", args[0])
	    }
	    passphrase, err := getPassphrase()
	    if err != nil {
	      return err
	    }
	    acc := account.Creation{
	      Name:       args[0],
	      Passphrase: passphrase,
	    }
	    fmt.Println("Enter mnemonic to recover keys from")
	    scanner := bufio.NewScanner(os.Stdin)
	    scanner.Scan()
	    m := scanner.Text()
	    if !bip39.IsMnemonicValid(m) {
	      return mnemonic.InvalidMnemonic
	    }
	    acc.Mnemonic = m
	    if err := account.CreateNewLocalAccount(&acc); err != nil {
	      return err
			}
			fmt.Println("Successfully recovered account from mnemonic!")
			return nil
		},
	}
  cmdRecoverMnemonic.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
  cmdRecoverMnemonic.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	cmdImportKS := &cobra.Command{
		Use:   "import-ks <KEYSTORE_FILE_PATH> [ACCOUNT_NAME]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Import an existing keystore key",
		RunE:  func(cmd *cobra.Command, args []string) error {
			userName := ""
			if len(args) == 2 {
				userName = args[1]
			}
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			name, err := account.ImportKeyStore(args[0], userName, passphrase)
			if !quietImport && err == nil {
				fmt.Printf("Imported keystore given account alias of `%s`\n", name)
			}
			return err
		},
	}
	cmdImportKS.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdImportKS.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")
	cmdImportKS.Flags().BoolVar(&quietImport, "quiet", false, "do not print out imported account name")

	cmdImportSK := &cobra.Command{
		Use:   "import-private-key <secp256k1_PRIVATE_KEY> [ACCOUNT_NAME]",
		Short: "Import an existing keystore key (only accept secp256k1 private keys)",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  func(cmd *cobra.Command, args []string) error {
			userName := ""
			if len(args) == 2 {
				userName = args[1]
			}
			passphrase, err := getPassphrase()
			if err != nil {
				return err
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
		Use:     "export-private-key <ACCOUNT_ADDRESS>",
		Short:   "Export the secp256k1 private key",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateAddress,
		RunE:    func(cmd *cobra.Command, args []string) error {
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			return account.ExportPrivateKey(addr.address, passphrase)
		},
	}
	cmdExportSK.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdExportSK.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	cmdExportKS := &cobra.Command{
		Use:     "export-ks <ACCOUNT_ADDRESS>",
		Short:   "Export the keystore file contents",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateAddress,
		RunE:    func(cmd *cobra.Command, args []string) error {
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			return account.ExportKeystore(addr.address, passphrase)
		},
	}
	cmdExportKS.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdExportKS.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	cmdGenerateBlsKey := &cobra.Command{
		Use:   "generate-bls-key",
		Short: "Generate bls keys then encrypt and save the private key with a requested passphrase",
		RunE:  func(cmd *cobra.Command, args []string) error {
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			return keys.GenBlsKeys(passphrase, blsFilePath)
		},
	}
	cmdGenerateBlsKey.Flags().StringVar(&blsFilePath, "bls-file-path", "",
		"absolute path of where to save encrypted bls private key")
	cmdGenerateBlsKey.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdGenerateBlsKey.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	cmdRecoverBlsKey := &cobra.Command{
		Use:   "recover-bls-key <ABSOLUTE_PATH_BLS_KEY>",
		Short: "Recover bls keys from an encrypted bls key file",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error {
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			return keys.RecoverBlsKeyFromFile(passphrase, args[0])
		},
	}
	cmdRecoverBlsKey.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdRecoverBlsKey.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	cmdSaveBlsKey := &cobra.Command{
		Use:   "save-bls-key <PRIVATE_BLS_KEY>",
		Short: "Encrypt and save the bls private key with a requested passphrase",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error {
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}
			return keys.SaveBlsKey(passphrase, blsFilePath, args[0])
		},
	}
	cmdSaveBlsKey.Flags().StringVar(&blsFilePath, "bls-file-path", "",
		"absolute path of where to save encrypted bls private key")
	cmdSaveBlsKey.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdSaveBlsKey.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	GetPublicBlsKey := &cobra.Command{
		Use:   "get-public-bls-key <PRIVATE_BLS_KEY>",
		Short: "Get the public bls key associated with the provided private bls key",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error {
			return keys.GetPublicBlsKey(args[0])
		},
	}

	return []*cobra.Command{cmdList, cmdLocation, cmdAdd, cmdRemove, cmdMnemonic, cmdRecoverMnemonic, cmdImportKS, cmdImportSK,
		cmdExportKS, cmdExportSK, cmdGenerateBlsKey, cmdRecoverBlsKey, cmdSaveBlsKey, GetPublicBlsKey}
}

func init() {
	cmdKeys := &cobra.Command{
		Use:   "keys",
		Short: "Add or view local private keys",
		Long:  "Manage your local keys",
		RunE:  func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdKeys.AddCommand(keysSub()...)
	RootCmd.AddCommand(cmdKeys)
}
