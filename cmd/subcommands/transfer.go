package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/go-sdk/pkg/validation"
	"github.com/harmony-one/harmony/accounts"

	"github.com/spf13/cobra"
)

var (
	fromAddress oneAddress
	toAddress   oneAddress
	amount      float64
	fromShardID uint32
	toShardID   uint32
	confirmWait uint32
	chainName   = chainIDWrapper{chainID: &common.Chain.TestNet}
	dryRun      bool
	unlockP     string
	gasPrice    int64
	setNonce    int64
)

func handlerForShard(senderShard uint32, node string) (*rpc.HTTPMessenger, error) {
	if checkNodeInput(node) {
		return rpc.NewHTTPHandler(node), nil
	}
	s, err := sharding.Structure(node)
	if err != nil {
		return nil, err
	}

	for _, shard := range s {
		if uint32(shard.ShardID) == senderShard {
			return rpc.NewHTTPHandler(shard.HTTP), nil
		}
	}

	return nil, nil
}

func opts(ctlr *transaction.Controller) {
	if dryRun {
		ctlr.Behavior.DryRun = true
	}
	if useLedgerWallet {
		ctlr.Behavior.SigningImpl = transaction.Ledger
	}
	if confirmWait > 0 {
		ctlr.Behavior.ConfirmationWaitTime = confirmWait
	}
}

func init() {
	cmdTransfer := &cobra.Command{
		Use:   "transfer",
		Short: "Create and send a transaction",
		Long: `
Create a transaction, sign it, and send off to the Harmony blockchain
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from := fromAddress.String()
			s, err := sharding.Structure(node)
			if err != nil {
				return err
			}
			err = validation.ValidShardIDs(fromShardID, toShardID, uint32(len(s)))
			if err != nil {
				return err
			}
			networkHandler, err := handlerForShard(fromShardID, node)
			if err != nil {
				return err
			}
			var ctrlr *transaction.Controller
			if useLedgerWallet {
				account := accounts.Account{Address: address.Parse(from)}
				ctrlr = transaction.NewController(networkHandler, nil, &account, *chainName.chainID, opts)
			} else {
				ks, acct, err := store.UnlockedKeystore(from, unlockP)
				if err != nil {
					return err
				}
				ctrlr = transaction.NewController(networkHandler, ks, acct, *chainName.chainID, opts)
			}

			if setNonce < 0 {
				setNonce = int64(getNextNonce(validatorAddress, networkHandler))
			}
			if transactionFailure := ctrlr.ExecuteTransaction(
				toAddress.String(),
				"",
				amount, gasPrice, setNonce,
				int(fromShardID),
				int(toShardID),
			); transactionFailure != nil {
				return transactionFailure
			}
			switch {
			case !dryRun && confirmWait == 0:
				fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, *ctrlr.ReceiptHash()))
			case !dryRun && confirmWait > 0:
				fmt.Println(common.ToJSONUnsafe(ctrlr.Receipt(), !noPrettyOutput))
			case dryRun:
				fmt.Println("Txn:")
				fmt.Println(ctrlr.TransactionToJSON(!noPrettyOutput))
				fmt.Println("RawTxn:", ctrlr.RawTransaction())
			}
			return nil
		},
	}

	cmdTransfer.Flags().Var(&fromAddress, "from", "sender's one address, keystore must exist locally")
	cmdTransfer.Flags().Var(&toAddress, "to", "the destination one address")
	cmdTransfer.Flags().BoolVar(&dryRun, "dry-run", false, "do not send signed transaction")
	cmdTransfer.Flags().Float64Var(&amount, "amount", 0.0, "amount")
	cmdTransfer.Flags().Int64Var(&gasPrice, "gas-price", 1, "gas price to pay")
	cmdTransfer.Flags().Int64Var(&setNonce, "nonce", -1, "set nonce for tx")
	cmdTransfer.Flags().Uint32Var(&fromShardID, "from-shard", 0, "source shard id")
	cmdTransfer.Flags().Uint32Var(&toShardID, "to-shard", 0, "target shard id")
	cmdTransfer.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	cmdTransfer.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	cmdTransfer.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock sender's keystore",
	)

	for _, flagName := range [...]string{"from", "to", "amount", "from-shard", "to-shard"} {
		cmdTransfer.MarkFlagRequired(flagName)
	}

	RootCmd.AddCommand(cmdTransfer)
}
