package v1

import (
	"fmt"

	rpcCommon "github.com/harmony-one/go-sdk/pkg/rpc/common"
)

const (
	prefix = "eth"
)

// Method is a list of known RPC methods
var Method = rpcCommon.RpcEnumList{
	GetShardingStructure:                    fmt.Sprintf("%s_getShardingStructure", prefix),
	GetNodeMetadata:                         fmt.Sprintf("%s_getNodeMetadata", prefix),
	GetLatestBlockHeader:                    fmt.Sprintf("%s_latestHeader", prefix),
	GetBlockByHash:                          fmt.Sprintf("%s_getBlockByHash", prefix),
	GetBlockByNumber:                        fmt.Sprintf("%s_getBlockByNumber", prefix),
	GetBlockTransactionCountByHash:          fmt.Sprintf("%s_getBlockTransactionCountByHash", prefix),
	GetBlockTransactionCountByNumber:        fmt.Sprintf("%s_getBlockTransactionCountByNumber", prefix),
	GetCode:                                 fmt.Sprintf("%s_getCode", prefix),
	GetTransactionByBlockHashAndIndex:       fmt.Sprintf("%s_getTransactionByBlockHashAndIndex", prefix),
	GetTransactionByBlockNumberAndIndex:     fmt.Sprintf("%s_getTransactionByBlockNumberAndIndex", prefix),
	GetTransactionByHash:                    fmt.Sprintf("%s_getTransactionByHash", prefix),
	GetStakingTransactionByHash:             fmt.Sprintf("%s_getStakingTransactionByHash", prefix),
	GetTransactionReceipt:                   fmt.Sprintf("%s_getTransactionReceipt", prefix),
	Syncing:                                 fmt.Sprintf("%s_syncing", prefix),
	PeerCount:                               "net_peerCount",
	GetBalance:                              fmt.Sprintf("%s_getBalance", prefix),
	GetStorageAt:                            fmt.Sprintf("%s_getStorageAt", prefix),
	GetTransactionCount:                     fmt.Sprintf("%s_getTransactionCount", prefix),
	SendTransaction:                         fmt.Sprintf("%s_sendTransaction", prefix),
	SendRawTransaction:                      fmt.Sprintf("%s_sendRawTransaction", prefix),
	Subscribe:                               fmt.Sprintf("%s_subscribe", prefix),
	GetPastLogs:                             fmt.Sprintf("%s_getLogs", prefix),
	GetWork:                                 fmt.Sprintf("%s_getWork", prefix),
	GetProof:                                fmt.Sprintf("%s_getProof", prefix),
	GetFilterChanges:                        fmt.Sprintf("%s_getFilterChanges", prefix),
	NewPendingTransactionFilter:             fmt.Sprintf("%s_newPendingTransactionFilter", prefix),
	NewBlockFilter:                          fmt.Sprintf("%s_newBlockFilter", prefix),
	NewFilter:                               fmt.Sprintf("%s_newFilter", prefix),
	Call:                                    fmt.Sprintf("%s_call", prefix),
	EstimateGas:                             fmt.Sprintf("%s_estimateGas", prefix),
	GasPrice:                                fmt.Sprintf("%s_gasPrice", prefix),
	BlockNumber:                             fmt.Sprintf("%s_blockNumber", prefix),
	UnSubscribe:                             fmt.Sprintf("%s_unsubscribe", prefix),
	NetVersion:                              "net_version",
	ProtocolVersion:                         fmt.Sprintf("%s_protocolVersion", prefix),
	SendRawStakingTransaction:               fmt.Sprintf("%s_sendRawStakingTransaction", prefix),
	GetElectedValidatorAddresses:            fmt.Sprintf("%s_getElectedValidatorAddresses", prefix),
	GetAllValidatorAddresses:                fmt.Sprintf("%s_getAllValidatorAddresses", prefix),
	GetValidatorInformation:                 fmt.Sprintf("%s_getValidatorInformation", prefix),
	GetAllValidatorInformation:              fmt.Sprintf("%s_getAllValidatorInformation", prefix),
	GetValidatorInformationByBlockNumber:    fmt.Sprintf("%s_getValidatorInformationByBlockNumber", prefix),
	GetAllValidatorInformationByBlockNumber: fmt.Sprintf("%s_getAllValidatorInformationByBlockNumber", prefix),
	GetDelegationsByDelegator:               fmt.Sprintf("%s_getDelegationsByDelegator", prefix),
	GetDelegationsByValidator:               fmt.Sprintf("%s_getDelegationsByValidator", prefix),
	GetCurrentTransactionErrorSink:          fmt.Sprintf("%s_getCurrentTransactionErrorSink", prefix),
	GetMedianRawStakeSnapshot:               fmt.Sprintf("%s_getMedianRawStakeSnapshot", prefix),
	GetCurrentStakingErrorSink:              fmt.Sprintf("%s_getCurrentStakingErrorSink", prefix),
	GetTransactionsHistory:                  fmt.Sprintf("%s_getTransactionsHistory", prefix),
	GetPendingTxnsInPool:                    fmt.Sprintf("%s_pendingTransactions", prefix),
	GetPendingCrosslinks:                    fmt.Sprintf("%s_getPendingCrossLinks", prefix),
	GetPendingCXReceipts:                    fmt.Sprintf("%s_getPendingCXReceipts", prefix),
	GetCurrentUtilityMetrics:                fmt.Sprintf("%s_getCurrentUtilityMetrics", prefix),
	ResendCX:                                fmt.Sprintf("%s_resendCx", prefix),
	GetSuperCommmittees:                     fmt.Sprintf("%s_getSuperCommittees", prefix),
	GetCurrentBadBlocks:                     fmt.Sprintf("%s_getCurrentBadBlocks", prefix),
	GetShardID:                              fmt.Sprintf("%s_getShardID", prefix),
	GetLastCrossLinks:                       fmt.Sprintf("%s_getLastCrossLinks", prefix),
	GetLatestChainHeaders:                   fmt.Sprintf("%s_getLatestChainHeaders", prefix),
}
