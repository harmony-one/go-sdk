package cmd

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	validatorSubCmds = []*cobra.Command{{
		Use:   "all-active",
		Short: "all validators marked as active",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetActiveValidatorAddresses, []interface{}{})
		},
	}, {
		Use:   "all",
		Short: "all validator addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetAllValidatorAddresses, []interface{}{})
		},
	}, {
		Use:   "information",
		Short: "original creation record of a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			addr := oneAddress{}
			if err := addr.Set(args[0]); err != nil {
				return err
			}
			return request(rpc.Method.GetValidatorInformation, []interface{}{args[0]})
		},
	},
	}
)
