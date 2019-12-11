package cmd

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

func init() {
	cmdFailures := &cobra.Command{
		Use:   "failures",
		Short: "Check in-memory record of failed transactions",
		Long:  `Check node for its in-memory record of failed transactions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
	cmdFailures.AddCommand([]*cobra.Command{{
		Use:   "staking",
		Short: "staking transaction failures",
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetCurrentTransactionErrorSink, []interface{}{})
		},
	}}...)
	RootCmd.AddCommand(cmdFailures)
}
