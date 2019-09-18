package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/harmony/accounts"

	"github.com/spf13/cobra"
)

var (
	fromAddress oneAddress
	toAddress   oneAddress
	amount      float64
	fromShardID int
	toShardID   int
	confirmWait uint32
	chainName   = chainIDWrapper{chainID: &common.Chain.TestNet}
	dryRun      bool
	unlockP     string
	gasPrice    float64
)

func handlerForShard(senderShard int, node string) (*rpc.HTTPMessenger, error) {
	s, err := sharding.Structure(node)
	if err != nil {
		return nil, err
	}
	for _, shard := range s {
		if shard.ShardID == senderShard {
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
			networkHandler, err := handlerForShard(fromShardID, node)
			if err != nil {
				return err
			}

			sender := address.Parse(from)
			var ctrlr *transaction.Controller
			if useLedgerWallet {
				account := accounts.Account{Address: sender}
				ctrlr = transaction.NewController(networkHandler, nil, &account, *chainName.chainID, opts)

			} else {
				ks := store.FromAddress(from)
				if ks == nil {
					return fmt.Errorf("could not open local keystore for %s", from)
				}

				account, lookupErr := ks.Find(accounts.Account{Address: sender})
				if lookupErr != nil {
					return fmt.Errorf("could not find %s in keystore", from)
				}
				if unlockError := ks.Unlock(account, unlockP); unlockError != nil {
					return unlockErr
				}
				ctrlr = transaction.NewController(networkHandler, ks, &account, *chainName.chainID, opts)
			}

			if transactionFailure := ctrlr.ExecuteTransaction(
				toAddress.String(),
				"",
				amount, gasPrice,
				fromShardID,
				toShardID,
			); transactionFailure != nil {
				return transactionFailure
			}
			switch {
			case !dryRun && confirmWait == 0:
				fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, *ctrlr.ReceiptHash()))
			case !dryRun && confirmWait > 0:
				fmt.Println(common.ToJSONUnsafe(ctrlr.Receipt(), !noPrettyOutput))
			case dryRun:
				fmt.Println(ctrlr.TransactionToJSON(!noPrettyOutput))
			}
			return nil
		},
	}

	cmdTransfer.Flags().Var(&fromAddress, "from", "sender's one address, keystore must exist locally")
	cmdTransfer.Flags().Var(&toAddress, "to", "the destination one address")
	cmdTransfer.Flags().BoolVar(&dryRun, "dry-run", false, "do not send signed transaction")
	cmdTransfer.Flags().Float64Var(&amount, "amount", 0.0, "amount")
	cmdTransfer.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	cmdTransfer.Flags().IntVar(&fromShardID, "from-shard", -1, "source shard id")
	cmdTransfer.Flags().IntVar(&toShardID, "to-shard", -1, "target shard id")
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
