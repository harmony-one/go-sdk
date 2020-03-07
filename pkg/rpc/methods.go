package rpc

import (
	"fmt"

	"github.com/pkg/errors"
)

// Using alias for now
type method = string

type errorCode int

type rpcEnumList struct {
	GetShardingStructure                method
	GetBlockByHash                      method
	GetBlockByNumber                    method
	GetBlockTransactionCountByHash      method
	GetBlockTransactionCountByNumber    method
	GetCode                             method
	GetTransactionByBlockHashAndIndex   method
	GetTransactionByBlockNumberAndIndex method
	GetTransactionByHash                method
	GetStakingTransactionByHash         method
	GetTransactionReceipt               method
	Syncing                             method
	PeerCount                           method
	GetBalance                          method
	GetStorageAt                        method
	GetTransactionCount                 method
	SendTransaction                     method
	SendRawTransaction                  method
	Subscribe                           method
	GetPastLogs                         method
	GetWork                             method
	GetProof                            method
	GetFilterChanges                    method
	NewPendingTransactionFilter         method
	NewBlockFilter                      method
	NewFilter                           method
	Call                                method
	EstimateGas                         method
	GasPrice                            method
	BlockNumber                         method
	UnSubscribe                         method
	NetVersion                          method
	ProtocolVersion                     method
	GetNodeMetadata                     method
	GetLatestBlockHeader                method
	SendRawStakingTransaction           method
	GetElectedValidatorAddresses        method
	GetAllValidatorAddresses            method
	GetValidatorInformation             method
	GetAllValidatorInformation          method
	GetDelegationsByDelegator           method
	GetDelegationsByValidator           method
	GetCurrentTransactionErrorSink      method
	GetValidatorMetrics                 method
	GetMedianRawStakeSnapshot           method
	GetCurrentStakingErrorSink          method
	GetTransactionsHistory              method
	GetPendingTxnsInPool                method
	GetPendingCrosslinks                method
	GetPendingCXReceipts                method
	GetCurrentUtilityMetrics            method
	ResendCX                            method
	GetSuperCommmittees                 method
	GetCurrentBadBlocks                 method
}

// Method is a list of known RPC methods
var Method = rpcEnumList{
	GetShardingStructure:                "hmy_getShardingStructure",
	GetNodeMetadata:                     "hmy_getNodeMetadata",
	GetLatestBlockHeader:                "hmy_latestHeader",
	GetBlockByHash:                      "hmy_getBlockByHash",
	GetBlockByNumber:                    "hmy_getBlockByNumber",
	GetBlockTransactionCountByHash:      "hmy_getBlockTransactionCountByHash",
	GetBlockTransactionCountByNumber:    "hmy_getBlockTransactionCountByNumber",
	GetCode:                             "hmy_getCode",
	GetTransactionByBlockHashAndIndex:   "hmy_getTransactionByBlockHashAndIndex",
	GetTransactionByBlockNumberAndIndex: "hmy_getTransactionByBlockNumberAndIndex",
	GetTransactionByHash:                "hmy_getTransactionByHash",
	GetStakingTransactionByHash:         "hmy_getStakingTransactionByHash",
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
	SendRawStakingTransaction:           "hmy_sendRawStakingTransaction",
	GetElectedValidatorAddresses:        "hmy_getElectedValidatorAddresses",
	GetAllValidatorAddresses:            "hmy_getAllValidatorAddresses",
	GetValidatorInformation:             "hmy_getValidatorInformation",
	GetAllValidatorInformation:          "hmy_getAllValidatorInformation",
	GetDelegationsByDelegator:           "hmy_getDelegationsByDelegator",
	GetDelegationsByValidator:           "hmy_getDelegationsByValidator",
	GetCurrentTransactionErrorSink:      "hmy_getCurrentTransactionErrorSink",
	GetValidatorMetrics:                 "hmy_getValidatorMetrics",
	GetMedianRawStakeSnapshot:           "hmy_getMedianRawStakeSnapshot",
	GetCurrentStakingErrorSink:          "hmy_getCurrentStakingErrorSink",
	GetTransactionsHistory:              "hmy_getTransactionsHistory",
	GetPendingTxnsInPool:                "hmy_pendingTransactions",
	GetPendingCrosslinks:                "hmy_getPendingCrossLinks",
	GetPendingCXReceipts:                "hmy_getPendingCXReceipts",
	GetCurrentUtilityMetrics:            "hmy_getCurrentUtilityMetrics",
	ResendCX:                            "hmy_resendCx",
	GetSuperCommmittees:                 "hmy_getSuperCommittees",
	GetCurrentBadBlocks:                 "hmy_getCurrentBadBlocks",
}

// TODO Use Reflection here to avoid typing out the cases

// ValidatedMethod checks if given method is known
func ValidatedMethod(m method) string {
	switch m := method(m); m {
	default:
		return string(m)
	}
}

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
