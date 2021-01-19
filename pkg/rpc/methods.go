package rpc

import (
	"fmt"

	rpcCommon "github.com/harmony-one/go-sdk/pkg/rpc/common"
	"github.com/pkg/errors"
)

var (
	RPCPrefix = "hmy"
	Method    rpcCommon.RpcEnumList
)

type RpcEnumList struct {
	GetShardingStructure                    rpcCommon.RpcMethod
	GetBlockByHash                          rpcCommon.RpcMethod
	GetBlockByNumber                        rpcCommon.RpcMethod
	GetBlockTransactionCountByHash          rpcCommon.RpcMethod
	GetBlockTransactionCountByNumber        rpcCommon.RpcMethod
	GetCode                                 rpcCommon.RpcMethod
	GetTransactionByBlockHashAndIndex       rpcCommon.RpcMethod
	GetTransactionByBlockNumberAndIndex     rpcCommon.RpcMethod
	GetTransactionByHash                    rpcCommon.RpcMethod
	GetStakingTransactionByHash             rpcCommon.RpcMethod
	GetTransactionReceipt                   rpcCommon.RpcMethod
	Syncing                                 rpcCommon.RpcMethod
	PeerCount                               rpcCommon.RpcMethod
	GetBalance                              rpcCommon.RpcMethod
	GetStorageAt                            rpcCommon.RpcMethod
	GetTransactionCount                     rpcCommon.RpcMethod
	SendTransaction                         rpcCommon.RpcMethod
	SendRawTransaction                      rpcCommon.RpcMethod
	Subscribe                               rpcCommon.RpcMethod
	GetPastLogs                             rpcCommon.RpcMethod
	GetWork                                 rpcCommon.RpcMethod
	GetProof                                rpcCommon.RpcMethod
	GetFilterChanges                        rpcCommon.RpcMethod
	NewPendingTransactionFilter             rpcCommon.RpcMethod
	NewBlockFilter                          rpcCommon.RpcMethod
	NewFilter                               rpcCommon.RpcMethod
	Call                                    rpcCommon.RpcMethod
	EstimateGas                             rpcCommon.RpcMethod
	GasPrice                                rpcCommon.RpcMethod
	BlockNumber                             rpcCommon.RpcMethod
	UnSubscribe                             rpcCommon.RpcMethod
	NetVersion                              rpcCommon.RpcMethod
	ProtocolVersion                         rpcCommon.RpcMethod
	GetNodeMetadata                         rpcCommon.RpcMethod
	GetLatestBlockHeader                    rpcCommon.RpcMethod
	SendRawStakingTransaction               rpcCommon.RpcMethod
	GetElectedValidatorAddresses            rpcCommon.RpcMethod
	GetAllValidatorAddresses                rpcCommon.RpcMethod
	GetValidatorInformation                 rpcCommon.RpcMethod
	GetAllValidatorInformation              rpcCommon.RpcMethod
	GetValidatorInformationByBlockNumber    rpcCommon.RpcMethod
	GetAllValidatorInformationByBlockNumber rpcCommon.RpcMethod
	GetDelegationsByDelegator               rpcCommon.RpcMethod
	GetDelegationsByValidator               rpcCommon.RpcMethod
	GetCurrentTransactionErrorSink          rpcCommon.RpcMethod
	GetMedianRawStakeSnapshot               rpcCommon.RpcMethod
	GetCurrentStakingErrorSink              rpcCommon.RpcMethod
	GetTransactionsHistory                  rpcCommon.RpcMethod
	GetPendingTxnsInPool                    rpcCommon.RpcMethod
	GetPendingCrosslinks                    rpcCommon.RpcMethod
	GetPendingCXReceipts                    rpcCommon.RpcMethod
	GetCurrentUtilityMetrics                rpcCommon.RpcMethod
	ResendCX                                rpcCommon.RpcMethod
	GetSuperCommmittees                     rpcCommon.RpcMethod
	GetCurrentBadBlocks                     rpcCommon.RpcMethod
	GetShardID                              rpcCommon.RpcMethod
	GetLastCrossLinks                       rpcCommon.RpcMethod
	GetLatestChainHeaders                   rpcCommon.RpcMethod
}

type errorCode int
type rpcErrorCodeList struct {
	rpcInvalidRequest       errorCode
	rpcMethodNotFound       errorCode
	rpcInvalidParams        errorCode
	rpcInternalError        errorCode
	rpcParseError           errorCode
	rpcMiscError            errorCode
	rpcTypeError            errorCode
	rpcInvalidAddressOrKey  errorCode
	rpcInvalidParameter     errorCode
	rpcDatabaseError        errorCode
	rpcDeserializationError errorCode
	rpcVerifyError          errorCode
	rpcVerifyRejected       errorCode
	rpcInWarmup             errorCode
	rpcMethodDeprecated     errorCode
	rpcGenericError         errorCode
}

// TODO Do not punt on the field names
var errorCodeEnumeration = rpcErrorCodeList{
	// Standard JSON-RPC 2.0 errors
	// RPC_INVALID_REQUEST is internally mapped to HTTP_BAD_REQUEST (400).
	// It should not be used for application-layer errors.
	rpcInvalidRequest: -32600,
	// RPC_METHOD_NOT_FOUND is internally mapped to HTTP_NOT_FOUND (404).
	// It should not be used for application-layer errors.
	rpcMethodNotFound: -32601,
	rpcInvalidParams:  -32602,
	// RPC_INTERNAL_ERROR should only be used for genuine errors in bitcoind
	// (for example datadir corruption).
	rpcInternalError: -32603,
	rpcParseError:    -32700,
	// General application defined errors
	rpcMiscError:            -1,     // std::exception thrown in command handling
	rpcTypeError:            -3,     // Unexpected type was passed as parameter
	rpcInvalidAddressOrKey:  -5,     // Invalid address or key
	rpcInvalidParameter:     -8,     // Invalid, missing or duplicate parameter
	rpcDatabaseError:        -20,    // Database error
	rpcDeserializationError: -22,    // Error parsing or validating structure in raw format
	rpcVerifyError:          -25,    // General error during transaction or block submission
	rpcVerifyRejected:       -26,    // Transaction or block was rejected by network rules
	rpcInWarmup:             -28,    // Client still warming up
	rpcMethodDeprecated:     -32,    // RPC method is deprecated
	rpcGenericError:         -32000, // Generic catchall error code
}

const (
	invalidRequestError      = "Invalid Request object"
	methodNotFoundError      = "Method not found"
	invalidParamsError       = "Invalid method parameter(s)"
	internalError            = "Internal JSON-RPC error"
	parseError               = "Error while parsing the JSON text"
	miscError                = "std::exception thrown in command handling"
	typeError                = "Unexpected type was passed as parameter"
	invalidAddressOrKeyError = "Invalid address or key"
	invalidParameterError    = "Invalid, missing or duplicate parameter"
	databaseError            = "Database error"
	deserializationError     = "Error parsing or validating structure in raw format"
	verifyError              = "General error during transaction or block submission"
	verifyRejectedError      = "Transaction or block was rejected by network rules"
	rpcInWarmupError         = "Client still warming up"
	methodDeprecatedError    = "RPC method deprecated"
	catchAllError            = "Catch all RPC error"
)

// ErrorCodeToError lifts an untyped error code from RPC to Error value
func ErrorCodeToError(message string, code float64) error {
	return errors.Wrap(errors.New(message), codeToMessage(code))
}

// TODO Use reflection here instead of typing out the cases or at least a map
func codeToMessage(err float64) string {
	switch e := errorCode(err); e {
	case errorCodeEnumeration.rpcInvalidRequest:
		return invalidRequestError
	case errorCodeEnumeration.rpcMethodNotFound:
		return methodNotFoundError
	case errorCodeEnumeration.rpcInvalidParams:
		return invalidParamsError
	case errorCodeEnumeration.rpcInternalError:
		return internalError
	case errorCodeEnumeration.rpcParseError:
		return parseError
	case errorCodeEnumeration.rpcMiscError:
		return miscError
	case errorCodeEnumeration.rpcTypeError:
		return typeError
	case errorCodeEnumeration.rpcInvalidAddressOrKey:
		return invalidAddressOrKeyError
	case errorCodeEnumeration.rpcInvalidParameter:
		return invalidParameterError
	case errorCodeEnumeration.rpcDatabaseError:
		return databaseError
	case errorCodeEnumeration.rpcDeserializationError:
		return deserializationError
	case errorCodeEnumeration.rpcVerifyError:
		return verifyError
	case errorCodeEnumeration.rpcVerifyRejected:
		return verifyRejectedError
	case errorCodeEnumeration.rpcInWarmup:
		return rpcInWarmupError
	case errorCodeEnumeration.rpcMethodDeprecated:
		return methodDeprecatedError
	case errorCodeEnumeration.rpcGenericError:
		return catchAllError
	default:
		return fmt.Sprintf("Error number %v not found", err)
	}
}
