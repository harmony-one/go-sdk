package cmd

import (
	"fmt"

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
		Use:     "metrics",
		Short:   "metrics about the performance of a validator",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateAddress,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			e := request(rpc.Method.GetValidatorMetrics, []interface{}{addr.address})
			if e != nil {
				fmt.Println("Metrics are only available for Validators that have participated in consensus committee.")
			}
			return e
		},
	}, {
		Use:     "information",
		Short:   "original creation record of a validator",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateAddress,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetValidatorInformation, []interface{}{addr.address})
		},
	}, {
		Use:     "all-information",
		Short:   "all validators information",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetAllValidatorInformation, []interface{}{args[0]})
		},
	},
	}
)
