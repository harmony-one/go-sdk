package cmd

import (
	"github.com/spf13/cobra"
)

// u, err := url.ParseRequestURI("http://google.com/")

var (
	cmdQuery = &cobra.Command{
		Use:   "account",
		Short: "Query account balance",
		Long:  `Query account balances`,
		Run: func(cmd *cobra.Command, args []string) {
			//
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
	// cmdQuery.AddCommand(&cobra.Command{
	// 	Use: ""
	// })
	RootCmd.AddCommand(cmdQuery)
}
