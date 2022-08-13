package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/governance"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/harmony/accounts"
	"github.com/spf13/cobra"
)

func init() {
	cmdGovernance := &cobra.Command{
		Use:   "governance",
		Short: "Interact with the Harmony spaces on https://snapshot.org",
		Long: `
Support interaction with the Harmony governance space on Snapshot, especially for validators that do not want to import their account private key into either metamask or onewallet.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdGovernance.AddCommand([]*cobra.Command{
		commandVote(),
	}...)

	RootCmd.AddCommand(cmdGovernance)
}

func commandVote() (cmd *cobra.Command) {
	var space string
	var proposal string
	var choice int
	var key string

	cmd = &cobra.Command{
		Use:   "vote-proposal",
		Short: "Vote on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyStore := store.FromAccountName(key)
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}

			if len(keyStore.Accounts()) <= 0 {
				return fmt.Errorf("Couldn't find address from the key")
			}

			account := accounts.Account{Address: keyStore.Accounts()[0].Address}
			err = keyStore.Unlock(accounts.Account{Address: keyStore.Accounts()[0].Address}, passphrase)
			if err != nil {
				return err
			}

			return governance.Vote(keyStore, account, space, proposal, choice)
		},
	}

	cmd.Flags().StringVar(&key, "key", "", "Account address. Must first use (hmy keys import-private-key) to import.")
	cmd.Flags().StringVar(&space, "space", "harmony-mainnet.eth", "Snapshot space")
	cmd.Flags().StringVar(&proposal, "proposal", "", "Proposal hash")
	// on snapshot they start with index 0 so use default value of -1 to make the choice explicit
	cmd.Flags().IntVar(&choice, "choice", -1, "Vote choice converted to integer")
	cmd.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("space")
	cmd.MarkFlagRequired("proposal")
	cmd.MarkFlagRequired("choose")
	return
}
