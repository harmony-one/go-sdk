package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

func init() {
	cmdBlockchain := &cobra.Command{
		Use:   "blockchain",
		Short: "Interact with the Harmony.one Blockchain",
		Long: `
Query Harmony's blockchain for high level metrics, queries
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	request := func(method rpc.RPCMethod, params interface{}) {
		fmt.Print(rpc.RPCRequest(method, node, params))
	}

	subCommands := [...]*cobra.Command{{
		Use:   "block-by-number",
		Short: "Get a harmony blockchain block by block number",
		Args:  cobra.ExactArgs(1),
		// TODO Add flag for second boolean parameter, consume argument
		Run: func(cmd *cobra.Command, args []string) {
			request(rpc.Method.GetBlockByNumber, []interface{}{args[0], true})
		},
	}, {
		Use:   "protocol-version",
		Short: "The version of the Harmony Protocol",
		Long: `
Query Harmony's blockchain for high level metrics, queries
`,
		Run: func(cmd *cobra.Command, args []string) {
			request(rpc.Method.ProtocolVersion, []string{})
		},
	}, {
		Use:   "transaction-by-hash",
		Short: "Get transaction by hash",
		Args:  cobra.ExactArgs(1),
		Long: `
Find a Harmony transaction by hash
`,
		Run: func(cmd *cobra.Command, args []string) {
			request(rpc.Method.GetTransactionByHash, []interface{}{args[0]})
		},
	}, {
		Use:   "transaction-by-receipt",
		Short: "Get transaction by receipt",
		Args:  cobra.ExactArgs(1),
		Long: `
Find a Harmony transaction by receipt
`,
		Run: func(cmd *cobra.Command, args []string) {
			request(rpc.Method.GetTransactionByHash, []interface{}{args[0]})
		},
	},
	}

	cmdBlockchain.Flags().StringVarP(
		&node,
		"node",
		"",
		"http://localhost:9500",
		"<host>:<port>",
	)

	cmdBlockchain.AddCommand(subCommands[:]...)

	RootCmd.AddCommand(cmdBlockchain)

}
