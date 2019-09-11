package cmd

import (
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	accountAddress string
)

func init() {
	cmdQuery := &cobra.Command{
		Use:   "account",
		Short: "Check account balance",
		Long:  `Query for the latest account balance given a Harmony Address`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			request(rpc.Method.GetBalance, []interface{}{address.ToBech32(address.Parse(args[0]))})
			// balanceRPCReply, failure := request(rpc.Method.GetBalance, []interface{}{address.ToBech32(address.Parse(args[0]))})
			// balance, _ := balanceRPCReply["result"].(string)
			// bln, _ := big.NewInt(0).SetString(balance[2:], 16)
			// fmt.Println(common.ConvertBalanceIntoReadableFormat(bln))
		},
	}

	RootCmd.AddCommand(cmdQuery)
}
