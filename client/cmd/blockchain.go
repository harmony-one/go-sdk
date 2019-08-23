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
			fmt.Print(rpc.RPCRequest(rpc.Method.ProtocolVersion, node, []string{}))
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
