package cmd

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	delegationSubCmds = []*cobra.Command{{
		Use:   "by-delegator",
		Short: "Print out the known chain-ids",
		Long:  "123",
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetDelegationsByDelegator, []interface{}{})
		},
	}, {
		Use:   "block-by-number",
		Short: "Get a harmony blockchain block by block number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetActiveValidatorAddresses, []interface{}{})
		},
	}}
)
