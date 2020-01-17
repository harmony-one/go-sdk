package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

func init() {
	cmdUtilities := &cobra.Command{
		Use:   "utility",
		Short: "common harmony blockchain utilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdUtilities.AddCommand(&cobra.Command{
		Use:   "metadata",
		Short: "data includes network specific values",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetNodeMetadata, []interface{}{})
		},
	})

	cmdMetrics := &cobra.Command{
		Use:   "metrics",
		Short: "mostly in-memory fluctuating values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdMetrics.AddCommand([]*cobra.Command{{
		Use:   "pending-crosslinks",
		Short: "dump the pending crosslinks in memory of target node",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetPendingCrosslinks, []interface{}{})
		},
	}, {
		Use:   "pending-cx-receipts",
		Short: "dump the pending cross shard receipts in memory of target node",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetPendingCXReceipts, []interface{}{})
		},
	},
	}...)

	cmdUtilities.AddCommand(cmdMetrics)

	cmdUtilities.AddCommand([]*cobra.Command{{
		Use:   "bech32-to-addr",
		Args:  cobra.ExactArgs(1),
		Short: "0x Address of a bech32 one-address",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := address.Bech32ToAddress(args[0])
			if err != nil {
				return err
			}
			fmt.Println(addr.Hex())
			return nil
		},
	}, {
		Use:   "addr-to-bech32",
		Args:  cobra.ExactArgs(1),
		Short: "bech32 one-address of a 0x Address",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO Implement this
			return nil
		},
	},
	}...)

	RootCmd.AddCommand(cmdUtilities)
}
