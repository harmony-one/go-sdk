package cmd

import (
	"fmt"
	"math/big"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/common/address"
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
		Run: func(cmd *cobra.Command, args []string) {
			// TODO Would we ever NOT want to have the latest?
			balanceRPCReply := rpc.RPCRequest(rpc.Method.GetBalance, node, []string{address.ToBech32(address.Parse(accountAddress)), "latest"})
			currentBalance, _ := balanceRPCReply["result"].(string)
			balance := big.NewInt(0)
			balance, _ = balance.SetString(currentBalance[2:], 16)
			fmt.Println(common.ConvertBalanceIntoReadableFormat(balance))
		},
	}

	cmdQuery.Flags().StringVar(&accountAddress, "address", "", "Harmony address of account to check")
	cmdQuery.MarkFlagRequired("address")
	RootCmd.AddCommand(cmdQuery)
}
