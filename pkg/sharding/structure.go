package sharding

import (
	"encoding/json"

	"github.com/harmony-one/go-sdk/pkg/rpc"
)

type RPCRoutes struct {
	HTTP    string `json:"http"`
	ShardID int    `json:"shardID"`
	WS      string `json:"ws"`
}

func Structure(node string) []RPCRoutes {
	type r struct {
		Result []RPCRoutes `json:"result"`
	}
	payload := rpc.RawRequest(rpc.Method.GetShardingStructure, node, []interface{}{})
	result := r{}
	json.Unmarshal(payload, &result)
	return result.Result
}
