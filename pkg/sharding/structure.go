package sharding

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
)

var (
	nanoAsDec = numeric.NewDec(denominations.Nano)
	oneAsDec  = numeric.NewDec(denominations.One)
)

// RPCRoutes reflects the RPC endpoints of the target network across shards
type RPCRoutes struct {
	HTTP    string `json:"http"`
	ShardID int    `json:"shardID"`
	WS      string `json:"ws"`
}

// Structure produces a slice of RPCRoutes for the network across shards
func Structure(node string) ([]RPCRoutes, error) {
	type r struct {
		Result []RPCRoutes `json:"result"`
	}
	p, e := rpc.RawRequest(rpc.Method.GetShardingStructure, node, []interface{}{})
	if e != nil {
		return nil, e
	}
	result := r{}
	if err := json.Unmarshal(p, &result); err != nil {
		return nil, err
	}
	return result.Result, nil
}

func CheckAllShards(node, oneAddr string, noPretty bool) (string, error) {
	var out bytes.Buffer
	out.WriteString("[")
	params := []interface{}{oneAddr, "latest"}
	s, err := Structure(node)
	if err != nil {
		return "", err
	}
	for i, shard := range s {
		balanceRPCReply, err := rpc.Request(rpc.Method.GetBalance, shard.HTTP, params)
		if err != nil {
			if common.DebugRPC {
				fmt.Printf("NOTE: Route %s failed.", shard.HTTP)
			}
			continue
		}
		if i != 0 {
			out.WriteString(",")
		}
		balance, _ := balanceRPCReply["result"].(string)
		bln := common.NewDecFromHex(balance)
		bln = bln.Quo(oneAsDec)
		out.WriteString(fmt.Sprintf(`{"shard":%d, "amount":%s}`,
			shard.ShardID,
			bln.String(),
		))
	}
	out.WriteString("]")
	if noPretty {
		return out.String(), nil
	}
	return common.JSONPrettyFormat(out.String()), nil
}
