package cmd

import (
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
			Short: "list all spaces of the governance app",
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
		Short: "list all proposals in one space",
		RunE: func(cmd *cobra.Command, args []string) error {
			return governance.PrintListProposals(space)
		},
	}

	cmd.Flags().StringVar(&space, "space", "", "list space")
	cmd.MarkFlagRequired("space")

	return
}

func commandViewProposal() (cmd *cobra.Command) {
	var proposal string

	cmd = &cobra.Command{
		Use:   "view-proposal",
		Short: "view one proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return governance.PrintViewProposal(proposal)
		},
	}

	cmd.Flags().StringVar(&proposal, "proposal", "", "proposal hash")
	cmd.MarkFlagRequired("proposal")

	return
}

func commandNewProposal() (cmd *cobra.Command) {
	var proposal string
	var key string

	cmd = &cobra.Command{
		Use:   "new-proposal",
		Short: "start a new proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyStore := store.FromAccountName(key)
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}

			account := accounts.Account{Address: keyStore.Accounts()[0].Address}
			err = keyStore.Unlock(accounts.Account{Address: keyStore.Accounts()[0].Address}, passphrase)
			if err != nil {
				return err
			}

			return governance.NewProposal(keyStore, account, proposal)
		},
	}

	cmd.Flags().StringVar(&proposal, "proposal-yaml", "", "proposal yaml path")
	cmd.Flags().StringVar(&key, "key", "", "private key name")
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
		Short: "vote one proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyStore := store.FromAccountName(key)
			passphrase, err := getPassphrase()
			if err != nil {
				return err
			}

			account := accounts.Account{Address: keyStore.Accounts()[0].Address}
			err = keyStore.Unlock(accounts.Account{Address: keyStore.Accounts()[0].Address}, passphrase)
			if err != nil {
				return err
			}

			return governance.Vote(keyStore, account, proposal, choice)
		},
	}

	cmd.Flags().StringVar(&proposal, "proposal", "", "proposal hash")
	cmd.Flags().StringVar(&choice, "choice", "", "choice")
	cmd.Flags().StringVar(&key, "key", "", "private key name")
	cmd.MarkFlagRequired("proposal")
	cmd.MarkFlagRequired("choose")
	cmd.MarkFlagRequired("key")
	return
}
