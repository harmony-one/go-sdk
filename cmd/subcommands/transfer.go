package cmd

import (
	"errors"
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

func handlerForShard(senderShard int, node string) *rpc.HTTPMessenger {
	for _, shard := range sharding.Structure(node) {
		if shard.ShardID == senderShard {
			return rpc.NewHTTPHandler(shard.HTTP)
		}
	}
	return nil
}

func init() {
	cmdTransfer := &cobra.Command{
		Use:   "transfer",
		Short: "Create and send a transaction",
		Long: `
Create a transaction, sign it, and send off to the Harmony blockchain
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler := handlerForShard(fromShardID, node)
			dryRunOpt := func(ctlr *transaction.Controller) {
				if dryRun {
					ctlr.Behavior.DryRun = true
				}
			}

			sender := address.Parse(fromAddress.String())
			var ctrlr *transaction.Controller
			if useLedgerWallet {
				account := accounts.Account{Address: sender}
				ctrlr = transaction.NewController(
					networkHandler, nil, &account,
					*chainName.chainID,
					dryRunOpt,
				)
				if transactionFailure := ctrlr.ExecuteHardwareTransaction(
					toAddress.String(),
					"",
					amount, gasPrice,
					fromShardID,
					toShardID,
				); transactionFailure != nil {
					return transactionFailure
				}

			} else {
				ks := store.FromAddress(fromAddress.String())
				if ks == nil {
					return fmt.Errorf("could not open local keystore for %s", fromAddress.String())
				}

				account, lookupErr := ks.Find(accounts.Account{Address: sender})
				if lookupErr != nil {
					return fmt.Errorf("could not find %s in keystore", fromAddress.String())
				}
				if unlockError := ks.Unlock(account, unlockP); unlockError != nil {
					return errors.New("could not unlock account with passphrase, perhaps need different phrase")
				}
				dryRunOpt := func(ctlr *transaction.Controller) {
					if dryRun {
						ctlr.Behavior.DryRun = true
					}
				}

				ctrlr = transaction.NewController(
					networkHandler, ks, &account,
					*chainName.chainID,
					dryRunOpt,
				)
				if transactionFailure := ctrlr.ExecuteTransaction(
					toAddress.String(),
					"",
					amount, gasPrice,
					fromShardID,
					toShardID,
				); transactionFailure != nil {
					return transactionFailure
				}

			}
			if !dryRun {
				fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, ctrlr.Receipt()))
			}
			return nil
		},
	}

	cmdTransfer.Flags().Var(&fromAddress, "from", "From can be an account alias or a one address")
	cmdTransfer.Flags().Var(&toAddress, "to", "the to address")
	cmdTransfer.Flags().BoolVar(&dryRun, "dry-run", false, "Do not send signed transaction")
	cmdTransfer.Flags().Float64Var(&amount, "amount", 0.0, "amount")
	cmdTransfer.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	cmdTransfer.Flags().IntVar(&fromShardID, "from-shard", -1, "source shard id")
	cmdTransfer.Flags().IntVar(&toShardID, "to-shard", -1, "target shard id")
	cmdTransfer.Flags().Var(&chainName, "chain-id", "What chain ID to target")
	cmdTransfer.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "Only waits if non-zero value, in seconds")
	cmdTransfer.Flags().StringVar(&unlockP, "passphrase", common.DefaultPassphrase, "Passphrase to unlock `from`")

	for _, flagName := range [...]string{"from", "to", "amount", "from-shard", "to-shard"} {
		cmdTransfer.MarkFlagRequired(flagName)
	}

	RootCmd.AddCommand(cmdTransfer)
}

// implemets pflag.Value interface
type oneAddress struct {
	address string
}

func (oneAddress oneAddress) String() string {
	return oneAddress.address
}

func (oneAddress *oneAddress) Set(s string) error {
	_, err := address.Bech32ToAddress(s)
	if err != nil {
		return err
	}
	oneAddress.address = s
	return nil
}

func (oneAddress oneAddress) Type() string {
	return "OneAddress"
}

// implements pflag.Value interface
type chainIDWrapper struct {
	chainID *common.ChainID
}

func (chainIDWrapper chainIDWrapper) String() string {
	return chainIDWrapper.chainID.Name
}

func (chainIDWrapper *chainIDWrapper) Set(s string) error {
	chainIDWrapper.chainID = common.StringToChainID(s)
	if chainIDWrapper.chainID == nil {
		return errors.New("ChainID \"" + s + "\" is invalid")
	}
	return nil
}

func (chainIDWrapper chainIDWrapper) Type() string {
	return "ChainID"
}
