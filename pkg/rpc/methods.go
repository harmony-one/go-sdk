package rpc

// Adapted from https://github.com/harmony-one/sdk/blob/master/packages/harmony-network/src/rpcMethod/rpc.ts
type RPCMethod string

type RPCErrorCode int

type rpcEnumList struct {
	GetShardingStructure                RPCMethod
	GetBlockByHash                      RPCMethod
	GetBlockByNumber                    RPCMethod
	GetBlockTransactionCountByHash      RPCMethod
	GetBlockTransactionCountByNumber    RPCMethod
	GetCode                             RPCMethod
	GetTransactionByBlockHashAndIndex   RPCMethod
	GetTransactionByBlockNumberAndIndex RPCMethod
	GetTransactionByHash                RPCMethod
	GetTransactionReceipt               RPCMethod
	Syncing                             RPCMethod
	PeerCount                           RPCMethod
	GetBalance                          RPCMethod
	GetStorageAt                        RPCMethod
	GetTransactionCount                 RPCMethod
	SendTransaction                     RPCMethod
	SendRawTransaction                  RPCMethod
	Subscribe                           RPCMethod
	GetPastLogs                         RPCMethod
	GetWork                             RPCMethod
	GetProof                            RPCMethod
	GetFilterChanges                    RPCMethod
	NewPendingTransactionFilter         RPCMethod
	NewBlockFilter                      RPCMethod
	NewFilter                           RPCMethod
	Call                                RPCMethod
	EstimateGas                         RPCMethod
	GasPrice                            RPCMethod
	BlockNumber                         RPCMethod
	UnSubscribe                         RPCMethod
	NetVersion                          RPCMethod
	ProtocolVersion                     RPCMethod
}

var Method = rpcEnumList{
	GetShardingStructure:                "hmy_getShardingStructure",
	GetBlockByHash:                      "hmy_getBlockByHash",
	GetBlockByNumber:                    "hmy_getBlockByNumber",
	GetBlockTransactionCountByHash:      "hmy_getBlockTransactionCountByHash",
	GetBlockTransactionCountByNumber:    "hmy_getBlockTransactionCountByNumber",
	GetCode:                             "hmy_getCode",
	GetTransactionByBlockHashAndIndex:   "hmy_getTransactionByBlockHashAndIndex",
	GetTransactionByBlockNumberAndIndex: "hmy_getTransactionByBlockNumberAndIndex",
	GetTransactionByHash:                "hmy_getTransactionByHash",
	GetTransactionReceipt:               "hmy_getTransactionReceipt",
	Syncing:                             "hmy_syncing",
	PeerCount:                           "net_peerCount",
	GetBalance:                          "hmy_getBalance",
	GetStorageAt:                        "hmy_getStorageAt",
	GetTransactionCount:                 "hmy_getTransactionCount",
	SendTransaction:                     "hmy_sendTransaction",
	SendRawTransaction:                  "hmy_sendRawTransaction",
	Subscribe:                           "hmy_subscribe",
	GetPastLogs:                         "hmy_getLogs",
	GetWork:                             "hmy_getWork",
	GetProof:                            "hmy_getProof",
	GetFilterChanges:                    "hmy_getFilterChanges",
	NewPendingTransactionFilter:         "hmy_newPendingTransactionFilter",
	NewBlockFilter:                      "hmy_newBlockFilter",
	NewFilter:                           "hmy_newFilter",
	Call:                                "hmy_call",
	EstimateGas:                         "hmy_estimateGas",
	GasPrice:                            "hmy_gasPrice",
	BlockNumber:                         "hmy_blockNumber",
	UnSubscribe:                         "hmy_unsubscribe",
	NetVersion:                          "net_version",
	ProtocolVersion:                     "hmy_protocolVersion",
}

type rpcErrorCodeList struct {
	RPC_INVALID_REQUEST        RPCErrorCode
	RPC_METHOD_NOT_FOUND       RPCErrorCode
	RPC_INVALID_PARAMS         RPCErrorCode
	RPC_INTERNAL_ERROR         RPCErrorCode
	RPC_PARSE_ERROR            RPCErrorCode
	RPC_MISC_ERROR             RPCErrorCode
	RPC_TYPE_ERROR             RPCErrorCode
	RPC_INVALID_ADDRESS_OR_KEY RPCErrorCode
	RPC_INVALID_PARAMETER      RPCErrorCode
	RPC_DATABASE_ERROR         RPCErrorCode
	RPC_DESERIALIZATION_ERROR  RPCErrorCode
	RPC_VERIFY_ERROR           RPCErrorCode
	RPC_VERIFY_REJECTED        RPCErrorCode
	RPC_IN_WARMUP              RPCErrorCode
	RPC_METHOD_DEPRECATED      RPCErrorCode
}

// TODO Turn these error codes into error values in query.go
var ErrorCode = rpcErrorCodeList{
	// Standard JSON-RPC 2.0 errors
	// RPC_INVALID_REQUEST is internally mapped to HTTP_BAD_REQUEST (400).
	// It should not be used for application-layer errors.
	RPC_INVALID_REQUEST: -32600,
	// RPC_METHOD_NOT_FOUND is internally mapped to HTTP_NOT_FOUND (404).
	// It should not be used for application-layer errors.
	RPC_METHOD_NOT_FOUND: -32601,
	RPC_INVALID_PARAMS:   -32602,
	// RPC_INTERNAL_ERROR should only be used for genuine errors in bitcoind
	// (for example datadir corruption).
	RPC_INTERNAL_ERROR: -32603,
	RPC_PARSE_ERROR:    -32700,
	// General application defined errors
	RPC_MISC_ERROR:             -1,  // std::exception thrown in command handling
	RPC_TYPE_ERROR:             -3,  // Unexpected type was passed as parameter
	RPC_INVALID_ADDRESS_OR_KEY: -5,  // Invalid address or key
	RPC_INVALID_PARAMETER:      -8,  // Invalid, missing or duplicate parameter
	RPC_DATABASE_ERROR:         -20, // Database error
	RPC_DESERIALIZATION_ERROR:  -22, // Error parsing or validating structure in raw format
	RPC_VERIFY_ERROR:           -25, // General error during transaction or block submission
	RPC_VERIFY_REJECTED:        -26, // Transaction or block was rejected by network rules
	RPC_IN_WARMUP:              -28, // Client still warming up
	RPC_METHOD_DEPRECATED:      -32, // RPC method is deprecated
}
