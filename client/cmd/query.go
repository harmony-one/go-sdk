package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

// u, err := url.ParseRequestURI("http://google.com/")

var (
	node     string
	cmdQuery = &cobra.Command{
		Use:   "account",
		Short: "Query account balance",
		Long:  `Query account balances`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(rpc.RPCRequest(node, "hmy_getBalance"))
		},
	}
)

func init() {
	cmdQuery.Flags().StringVarP(
		&node,
		"node",
		"",
		"http://localhost:9500",
		"<host>:<port>",
	)
	RootCmd.AddCommand(cmdQuery)
}
