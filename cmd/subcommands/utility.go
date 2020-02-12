package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/shard"
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
		Use:   "shard-for-bls",
		Args:  cobra.ExactArgs(1),
		Short: "which shard (default assumes mainnet) this BLS key would be assigned to",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := bls.PublicKey{}
			if err := key.DeserializeHexStr(args[0]); err != nil {
				return err
			}
			// TODO Need to take flag changing the shard count per chainID
			shardBig := big.NewInt(4)
			wrapper := shard.FromLibBLSPublicKeyUnsafe(&key)
			shardID := int(new(big.Int).Mod(wrapper.Big(), shardBig).Int64())
			type t struct {
				ShardID int `json:"shard-id"`
			}
			result, err := json.Marshal(t{shardID})
			if err != nil {
				return err
			}

			fmt.Println(string(result))
			return nil
		},
	},
	}...)

	RootCmd.AddCommand(cmdUtilities)
}
