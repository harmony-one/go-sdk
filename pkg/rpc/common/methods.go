package common

type RpcMethod = string

type RpcEnumList struct {
	GetShardingStructure                    RpcMethod
	GetBlockByHash                          RpcMethod
	GetBlockByNumber                        RpcMethod
	GetBlockTransactionCountByHash          RpcMethod
	GetBlockTransactionCountByNumber        RpcMethod
	GetCode                                 RpcMethod
	GetTransactionByBlockHashAndIndex       RpcMethod
	GetTransactionByBlockNumberAndIndex     RpcMethod
	GetTransactionByHash                    RpcMethod
	GetStakingTransactionByHash             RpcMethod
	GetTransactionReceipt                   RpcMethod
	Syncing                                 RpcMethod
	PeerCount                               RpcMethod
	GetBalance                              RpcMethod
	GetStorageAt                            RpcMethod
	GetTransactionCount                     RpcMethod
	SendTransaction                         RpcMethod
	SendRawTransaction                      RpcMethod
	Subscribe                               RpcMethod
	GetPastLogs                             RpcMethod
	GetWork                                 RpcMethod
	GetProof                                RpcMethod
	GetFilterChanges                        RpcMethod
	NewPendingTransactionFilter             RpcMethod
	NewBlockFilter                          RpcMethod
	NewFilter                               RpcMethod
	Call                                    RpcMethod
	EstimateGas                             RpcMethod
	GasPrice                                RpcMethod
	BlockNumber                             RpcMethod
	UnSubscribe                             RpcMethod
	NetVersion                              RpcMethod
	ProtocolVersion                         RpcMethod
	GetNodeMetadata                         RpcMethod
	GetLatestBlockHeader                    RpcMethod
	SendRawStakingTransaction               RpcMethod
	GetElectedValidatorAddresses            RpcMethod
	GetAllValidatorAddresses                RpcMethod
	GetValidatorInformation                 RpcMethod
	GetAllValidatorInformation              RpcMethod
	GetValidatorInformationByBlockNumber    RpcMethod
	GetAllValidatorInformationByBlockNumber RpcMethod
	GetDelegationsByDelegator               RpcMethod
	GetDelegationsByValidator               RpcMethod
	GetCurrentTransactionErrorSink          RpcMethod
	GetMedianRawStakeSnapshot               RpcMethod
	GetCurrentStakingErrorSink              RpcMethod
	GetTransactionsHistory                  RpcMethod
	GetPendingTxnsInPool                    RpcMethod
	GetPendingCrosslinks                    RpcMethod
	GetPendingCXReceipts                    RpcMethod
	GetCurrentUtilityMetrics                RpcMethod
	ResendCX                                RpcMethod
	GetSuperCommmittees                     RpcMethod
	GetCurrentBadBlocks                     RpcMethod
	GetShardID                              RpcMethod
	GetLastCrossLinks                       RpcMethod
	GetLatestChainHeaders                   RpcMethod
}

// ValidatedMethod checks if given method is known
func ValidatedMethod(m RpcMethod) string {
	switch m := RpcMethod(m); m {
	default:
		return string(m)
	}
}
