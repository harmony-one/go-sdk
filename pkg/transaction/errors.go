package transaction

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/harmony-one/go-sdk/pkg/rpc"
)

var (
	errorSinkRPCs = []string{
		rpc.Method.GetCurrentTransactionErrorSink,
		rpc.Method.GetCurrentStakingErrorSink,
	}
)

// Error is the struct for all errors in the error sinks.
// Note that StakingDirective is non-nil for staking txn errors.
type Error struct {
	TxHashID             *string `json:"tx-hash-id"`
	StakingDirective     *string `json:"directive-kind"` // optional error field for staking
	ErrMessage           *string `json:"error-message"`
	TimestampOfRejection int64   `json:"time-at-rejection"`
}

// Error an error with the ErrMessage and TimestampOfRejection as the message.
func (e *Error) Error() error {
	if e.StakingDirective != nil {
		return fmt.Errorf("[%s] %s -- %s",
			*e.StakingDirective, *e.ErrMessage, time.Unix(e.TimestampOfRejection, 0).String(),
		)
	} else {
		return fmt.Errorf("[Plain transaction] %s -- %s",
			*e.ErrMessage, time.Unix(e.TimestampOfRejection, 0).String(),
		)
	}
}

// Errors ...
type Errors []*Error

func getTxErrorBySink(txHash, errorSinkRPC string, messenger rpc.T) (Errors, error) {
	var txErrors Errors
	response, err := messenger.SendRPC(errorSinkRPC, []interface{}{})
	if err != nil {
		return nil, err
	}
	var allErrors []Error
	asJSON, _ := json.Marshal(response["result"])
	_ = json.Unmarshal(asJSON, &allErrors)
	for i := range allErrors {
		txError := allErrors[i]
		if *txError.TxHashID == txHash {
			txErrors = append(txErrors, &txError)
		}
	}
	return txErrors, nil
}

// GetError returns all errors for a given (staking or plain) transaction hash.
func GetError(txHash string, messenger rpc.T) (Errors, error) {
	for _, errorSinkRpc := range errorSinkRPCs {
		txErrors, err := getTxErrorBySink(txHash, errorSinkRpc, messenger)
		if err != nil {
			return Errors{}, err
		}
		if txErrors != nil {
			return txErrors, nil
		}
	}
	return Errors{}, nil
}
