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
		Use:   "plain",
		Short: "plain transaction failures",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetCurrentTransactionErrorSink, []interface{}{})
		},
	}, {
		Use:   "staking",
		Short: "staking transaction failures",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetCurrentStakingErrorSink, []interface{}{})
		},
	},
	}...)
	RootCmd.AddCommand(cmdFailures)
}
