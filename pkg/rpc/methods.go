package rpc

import (
	"github.com/pkg/errors"
	"fmt"
)

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
	rpcInvalidRequest        RPCErrorCode
	rpcMethodNotFound        RPCErrorCode
	rpcInvalidParams         RPCErrorCode
	rpcInternalError         RPCErrorCode
	rpcParseError            RPCErrorCode
	rpcMiscError             RPCErrorCode
	rpcTypeError             RPCErrorCode
	rpcInvalidAddressOrKey   RPCErrorCode
	rpcInvalidParameter      RPCErrorCode
	rpcDatabaseError         RPCErrorCode
	rpcDeserializationError  RPCErrorCode
	rpcVerifyError           RPCErrorCode
	rpcVerifyRejected        RPCErrorCode
	rpcInWarmup              RPCErrorCode
	rpcMethodDeprecated      RPCErrorCode
}

// TODO Turn these error codes into error values in query.go
var ErrorCode = rpcErrorCodeList{
	// Standard JSON-RPC 2.0 errors
	// RPC_INVALID_REQUEST is internally mapped to HTTP_BAD_REQUEST (400).
	// It should not be used for application-layer errors.
	rpcInvalidRequest: -32600,
	// RPC_METHOD_NOT_FOUND is internally mapped to HTTP_NOT_FOUND (404).
	// It should not be used for application-layer errors.
	rpcMethodNotFound: -32601,
	rpcInvalidParams:   -32602,
	// RPC_INTERNAL_ERROR should only be used for genuine errors in bitcoind
	// (for example datadir corruption).
	rpcInternalError: -32603,
	rpcParseError:    -32700,
	// General application defined errors
	rpcMiscError:              -1, // std::exception thrown in command handling
	rpcTypeError:              -3, // Unexpected type was passed as parameter
	rpcInvalidAddressOrKey:    -5, // Invalid address or key
	rpcInvalidParameter:       -8, // Invalid, missing or duplicate parameter
	rpcDatabaseError:         -20, // Database error
	rpcDeserializationError:  -22, // Error parsing or validating structure in raw format
	rpcVerifyError:           -25, // General error during transaction or block submission
	rpcVerifyRejected:        -26, // Transaction or block was rejected by network rules
	rpcInWarmup:              -28, // Client still warming up
	rpcMethodDeprecated:      -32, // RPC method is deprecated
}

const (
	invalidRequestError       = "Invalid Request object"
	methodNotFoundError       = "Method not found"
	invalidParamsError        = "Invalid method parameter(s)"
	internalError             = "Internal JSON-RPC error"
	parseError                = "Error while parsing the JSON text"
	miscError                 = "std::exception thrown in command handling"
	typeError                 = "Unexpected type was passed as parameter"
	invalidAddressOrKeyError  = "Invalid address or key"
	invalidParameterError     = "Invalid, missing or duplicate parameter"
	databaseError             = "Database error"
	deserializationError      = "Error parsing or validating structure in raw format"
	verifyError               = "General error during transaction or block submission"
	verifyRejectedError       = "Transaction or block was rejected by network rules"
	rpcInWarmupError          = "Client still warming up"
	methodDeprecatedError     = "RPC method is deprecated"
)

func ErrorNumberToError(message string,err float64) error {
	return errors.Wrap(errors.New(message),getErrorMessage(err))
}

func getErrorMessage(err float64) string {
	errorcode := RPCErrorCode(err)
	switch errorcode {
	case ErrorCode.rpcInvalidRequest:
		return invalidRequestError
	case ErrorCode.rpcMethodNotFound:
		return methodNotFoundError
	case ErrorCode.rpcInvalidParams:
		return invalidParamsError
	case ErrorCode.rpcInternalError:
		return internalError
	case ErrorCode.rpcParseError:
		return parseError
	case ErrorCode.rpcMiscError:
		return miscError
	case ErrorCode.rpcTypeError:
		return typeError
	case ErrorCode.rpcInvalidAddressOrKey:
		return invalidAddressOrKeyError
	case ErrorCode.rpcInvalidParameter:
		return invalidParameterError
	case ErrorCode.rpcDatabaseError:
		return databaseError
	case ErrorCode.rpcDeserializationError:
		return deserializationError
	case ErrorCode.rpcVerifyError:
		return verifyError
	case ErrorCode.rpcVerifyRejected:
		return verifyRejectedError
	case ErrorCode.rpcInWarmup:
		return rpcInWarmupError
	case ErrorCode.rpcMethodDeprecated:
		return methodDeprecatedError
	default:
		panic(fmt.Sprintf("Error number %v not found", err))
	}
}