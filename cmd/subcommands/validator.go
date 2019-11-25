package cmd

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	validatorSubCmds = []*cobra.Command{{
		Use:   "all-active",
		Short: "Get a harmony blockchain block by block number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetActiveValidatorAddresses, []interface{}{})
		},
	}, {
		Use:   "all",
		Short: "Print out the known chain-ids",
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetAllValidatorAddresses, []interface{}{})
		},
	}, {
		Use:   "information",
		Short: "Get a harmony blockchain block by block number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetActiveValidatorAddresses, []interface{}{})

		},
	},
	}
)
