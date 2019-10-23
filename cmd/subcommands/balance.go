package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/spf13/cobra"
)

func init() {
	cmdQuery := &cobra.Command{
		Use:   "balances",
		Short: "Check account balance on all shards",
		Long:  `Query for the latest account balance given a Harmony Address`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := sharding.CheckAllShards(node, args[0], noPrettyOutput)
			if err != nil {
				return err
			}
			fmt.Println(r)
			return nil
		},
	}

	RootCmd.AddCommand(cmdQuery)
}
