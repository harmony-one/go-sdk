package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

func init() {
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
	},
	}
	cmdBlockchain.AddCommand(subCommands[:]...)
	RootCmd.AddCommand(cmdBlockchain)
}
