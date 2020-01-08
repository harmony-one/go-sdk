package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/spf13/cobra"
)

func init() {
	cmdUtilities := &cobra.Command{
		Use:   "utility",
		Short: "Check in-memory record of failed transactions",
		Long:  `Check node for its in-memory record of failed transactions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
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
