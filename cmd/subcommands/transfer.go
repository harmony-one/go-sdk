package cmd

import (
	"fmt"
	"os"

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
	fromAddress string
	toAddress   string
	amount      float64
	fromShardID int
	toShardID   int
	confirmWait uint32
	accountName string
	chainName   string
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
		Run: func(cmd *cobra.Command, args []string) {
			networkHandler := handlerForShard(fromShardID, node)
			ks := store.FromAccountName(accountName)
			sender := address.Parse(fromAddress)
			account, _ := ks.Find(accounts.Account{Address: sender})
			ks.Unlock(account, common.DefaultPassphrase)
			fromCmdLineFlags := func(ctlr *transaction.Controller) {
				//
			}

			ctrlr, err := transaction.NewController(
				networkHandler, ks, &account,
				*common.StringToChainID(chainName),
				fromCmdLineFlags,
			)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			if transactionFailure := ctrlr.ExecuteTransaction(
				toAddress,
				"",
				amount,
				fromShardID,
				toShardID,
			); transactionFailure != nil {
				fmt.Println(transactionFailure)
				os.Exit(-1)
			}
			fmt.Println(ctrlr.Receipt())
		},
	}

	// TODO Intern do custom flag validation for one address: see https://github.com/spf13/cobra/issues/376
	cmdTransfer.Flags().StringVar(&fromAddress, "from", "", "From can be an account alias or a one address")
	cmdTransfer.Flags().StringVar(&toAddress, "to", "", "the to address")

	cmdTransfer.Flags().StringVar(&accountName, "account-name", "", "account-name")

	cmdTransfer.Flags().Float64Var(&amount, "amount", 0.0, "amount")
	cmdTransfer.Flags().IntVar(&fromShardID, "from-shard", -1, "source shard id")
	cmdTransfer.Flags().IntVar(&toShardID, "to-shard", -1, "target shard id")
	cmdTransfer.Flags().StringVar(&chainName, "chain-id", common.Chain.TestNet.Name, "What chain ID to target")
	cmdTransfer.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "Only waits if non-zero value, in seconds")

	for _, flagName := range [...]string{"from", "to", "amount", "from-shard", "to-shard"} {
		cmdTransfer.MarkFlagRequired(flagName)
	}

	RootCmd.AddCommand(cmdTransfer)
}
