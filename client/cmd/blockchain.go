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

	subCommands := [...]*cobra.Command{{
		Use:   "update",
		Short: "Change the password used to protect private key",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd)
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
