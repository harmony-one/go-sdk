package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

func init() {
	cmdValidator := &cobra.Command{
		Use:   "validator",
		Short: "information about validators",
		Long: `
Look up information about validator information
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmdDelegation := &cobra.Command{
		Use:   "delegation",
		Short: "information about delegations",
		Long: `
Look up information about delegation
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmdBlockchain := &cobra.Command{
		Use:   "blockchain",
		Short: "Interact with the Harmony.one Blockchain",
		Long: `
Query Harmony's blockchain for completed transaction, historic records
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	subCommands := []*cobra.Command{{
		Use:   "block-by-number",
		Short: "Get a harmony blockchain block by block number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetBlockByNumber, []interface{}{args[0], true})
		},
	}, {
		Use:   "known-chains",
		Short: "Print out the known chain-ids",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(common.ToJSONUnsafe(common.AllChainIDs(), !noPrettyOutput))
		},
	}, {
		Use:   "protocol-version",
		Short: "The version of the Harmony Protocol",
		Long: `
Query Harmony's blockchain for high level metrics, queries
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.ProtocolVersion, []interface{}{})
		},
	}, {
		Use:   "transaction-by-hash",
		Short: "Get transaction by hash",
		Args:  cobra.ExactArgs(1),
		Long: `
Inputs of a transaction and r, s, v value of transaction
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetTransactionByHash, []interface{}{args[0]})
		},
	}, {
		Use:   "transaction-receipt",
		Short: "Get information about a finalized transaction",
		Args:  cobra.ExactArgs(1),
		Long: `
High level information about transaction, like blockNumber, blockHash
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetTransactionReceipt, []interface{}{args[0]})
		},
	}, {
		Use:   "current-nonce",
		Short: "Current nonce of an account",
		Args:  cobra.ExactArgs(1),
		Long:  `Current nonce number of a one-address`,
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := oneAddress{}
			if err := addr.Set(args[0]); err != nil {
				return err
			}
			return request(rpc.Method.GetTransactionCount, []interface{}{args[0]})
		},
	}, {
		Use:   "latest-header",
		Short: "Get the latest header",
		RunE: func(cmd *cobra.Command, arg []string) error {
			noLatest = true
			return request(rpc.Method.GetLatestBlockHeader, []interface{}{})
		},
	},
	}

	cmdBlockchain.AddCommand(cmdValidator)
	cmdBlockchain.AddCommand(cmdDelegation)
	cmdValidator.AddCommand(validatorSubCmds[:]...)
	cmdDelegation.AddCommand(delegationSubCmds[:]...)
	cmdBlockchain.AddCommand(subCommands[:]...)
	RootCmd.AddCommand(cmdBlockchain)
}
