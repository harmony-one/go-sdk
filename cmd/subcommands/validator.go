package cmd

import (
	"fmt"
	"strconv"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	validatorSubCmds = []*cobra.Command{{
		Use:   "elected",
		Short: "elected validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetElectedValidatorAddresses, []interface{}{})
		},
	}, {
		Use:   "all",
		Short: "all validator addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.GetAllValidatorAddresses, []interface{}{})
		},
	}, {
		Use:     "information",
		Short:   "original creation record of a validator",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateAddress,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			e := request(rpc.Method.GetValidatorInformation, []interface{}{addr.address})
			if e != nil {
				return fmt.Errorf("validator address not found: %s", addr.address)
			}
			return e
		},
	}, {
		Use:     "information-by-block-number",
		Short:   "original creation record of a validator by block number",
		Args:    cobra.ExactArgs(2),
		PreRunE: validateAddress,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetValidatorInformationByBlockNumber, []interface{}{addr.address, args[1]})
		},
	}, {
		Use:   "all-information",
		Short: "all validators information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			page, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "the page argument must be integer, supplied %v", args[0])
			}
			return request(rpc.Method.GetAllValidatorInformation, []interface{}{page})
		},
	}, {
		Use:   "all-information-by-block-number",
		Short: "all validators information by block number",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			page, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "the page argument must be integer, supplied %v", args[0])
			}
			return request(rpc.Method.GetAllValidatorInformationByBlockNumber, []interface{}{page, args[1]})
		},
	},
	}
)
