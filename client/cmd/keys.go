package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/keys"
	"github.com/spf13/cobra"
)

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

	cmdMnemonic := &cobra.Command{
		Use:   "mnemonic",
		Short: "Compute the bip39 mnemonic for some input entropy",
		Run: func(cmd *cobra.Command, args []string) {
			mnemonic := keys.GenerateMnemonic()
			fmt.Println(mnemonic)
			fmt.Println(keys.NewAccountByMnemonic(mnemonic))
		},
	}


	var password string
	cmdAdd := &cobra.Command{
		Use:   "add",
		Short: "Create a new key with passphrase",
		Run: func(cmd *cobra.Command, args []string) {
			keys.AddNewKey(password)
		},
	}

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		Run: func(cmd *cobra.Command, args []string) {
			keys.ListKeys(keyStoreDir)
		},
	}

	cmdShow := &cobra.Command{
		Use:   "show",
		Short: "Show key info for the given name",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd)
		},
	}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: "Delete the given key",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd)
		},
	}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: "Change the password used to protect private key",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd)
		},
	}

	cmdExport := &cobra.Command{
		Use:   "export",
		Short: "Export your keystore",
		Run: func(cmd *cobra.Command, args []string) {
			// keys.ListKeys(keyStoreDir)
		},
	}

	cmdAdd.Flags().StringVarP(&password, "password", "w", "", "password to encrypt generated private key")
	cmdAdd.MarkFlagRequired("password")

	cmdKeys.AddCommand(cmdMnemonic, cmdAdd, cmdList, cmdShow, cmdDelete, cmdUpdate, cmdExport)
	RootCmd.AddCommand(cmdKeys)

}
