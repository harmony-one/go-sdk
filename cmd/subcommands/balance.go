package cmd

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/spf13/cobra"
)

func checkAllShards(node, addr string, noPretty bool) string {
	var out bytes.Buffer
	out.WriteString("[")
	params := []interface{}{addr, "latest"}
	s := sharding.Structure(node)
	for i, shard := range s {
		balanceRPCReply, _ := rpc.Request(rpc.Method.GetBalance, shard.HTTP, params)
		balance, _ := balanceRPCReply["result"].(string)
		bln, _ := big.NewInt(0).SetString(balance[2:], 16)
		out.WriteString(fmt.Sprintf(`{"shard":%d, "amount":%s}`,
			shard.ShardID,
			common.ConvertBalanceIntoReadableFormat(bln),
		))
		if i != len(s)-1 {
			out.WriteString(",")
		}
	}
	out.WriteString("]")
	if noPretty {
		return out.String()
	}
	return common.JSONPrettyFormat(out.String())
}

func init() {
	cmdQuery := &cobra.Command{
		Use:   "balance",
		Short: "Check account balance",
		Long:  `Query for the latest account balance given a Harmony Address`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(checkAllShards(node, args[0], noPrettyOutput))
		},
	}

	RootCmd.AddCommand(cmdQuery)
}
