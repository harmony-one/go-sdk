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
		Short: "Support interaction with the Harmony governance app.",
		Long: `
Support interaction with the Harmony governance app, especially for validators that do not want to import their account private key into either metamask or onewallet.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdGovernance.AddCommand([]*cobra.Command{
		{
			Use:   "list-space",
			Short: "List all spaces of the governance app",
			RunE: func(cmd *cobra.Command, args []string) error {
				return governance.PrintListSpace()
			},
		},
		commandListProposal(),
		commandViewProposal(),
		commandNewProposal(),
		commandVote(),
	}...)

	RootCmd.AddCommand(cmdGovernance)
}

func commandListProposal() (cmd *cobra.Command) {
	var space string

	cmd = &cobra.Command{
		Use:   "list-proposal",
		Short: "List all proposals for the given space",
		RunE: func(cmd *cobra.Command, args []string) error {
			return governance.PrintListProposals(space)
		},
	}

	cmd.Flags().StringVar(&space, "space", "", "Space the proposal belongs to e.g. 'staking-mainnet'")
	cmd.MarkFlagRequired("space")

	return
}

func commandViewProposal() (cmd *cobra.Command) {
	var proposal string

	cmd = &cobra.Command{
		Use:   "view-proposal",
		Short: "View a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return governance.PrintViewProposal(proposal)
		},
	}

	cmd.Flags().StringVar(&proposal, "proposal", "", "Proposal hash")
	cmd.MarkFlagRequired("proposal")

	return
}

func commandNewProposal() (cmd *cobra.Command) {
	var proposal string
	var key string

	cmd = &cobra.Command{
		Use:   "new-proposal",
		Short: "Start a new proposal",
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

			return governance.NewProposal(keyStore, account, proposal)
		},
	}

	cmd.Flags().StringVar(&proposal, "proposal-yaml", "", "Proposal yaml path")
	cmd.Flags().StringVar(&key, "key", "", "Account address. Must first use (hmy keys import-private-key) to import.")
	cmd.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmd.MarkFlagRequired("proposal-yaml")
	cmd.MarkFlagRequired("key")

	return
}

func commandVote() (cmd *cobra.Command) {
	var proposal string
	var choice string
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

			return governance.Vote(keyStore, account, proposal, choice)
		},
	}

	cmd.Flags().StringVar(&proposal, "proposal", "", "Proposal hash")
	cmd.Flags().StringVar(&choice, "choice", "", "Vote choice e.g. 'agree' or 'disagree'")
	cmd.Flags().StringVar(&key, "key", "", "Account address. Must first use (hmy keys import-private-key) to import.")
	cmd.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmd.MarkFlagRequired("proposal")
	cmd.MarkFlagRequired("choose")
	cmd.MarkFlagRequired("key")
	return
}
