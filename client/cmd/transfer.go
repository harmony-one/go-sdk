package cmd

import (
	"fmt"
	"os"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/harmony/accounts"

	"github.com/harmony-one/go-sdk/pkg/rpc"
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
	chainID     uint32
)

func init() {
	cmdTransfer := &cobra.Command{
		Use:   "transfer",
		Short: "Create and send a transaction",
		Long: `
Create a transaction, sign it, and send off to the Harmony blockchain
`,
		Run: func(cmd *cobra.Command, args []string) {
			networkHandler := rpc.NewHTTPHandler(node)
			ks := store.FromAccountName(accountName)
			sender := address.Parse(fromAddress)
			account, _ := ks.Find(accounts.Account{Address: sender})
			// _ :=
			ks.Unlock(account, common.DefaultPassphrase)
			// ks.SignTx(a accounts.Account, tx *types.Transaction, chainID *big.Int)
			fromCmdLineFlags := func(ctlr *transaction.Controller) {
				// ctlr.PreferOneAddress = true
				// if confirmWait != 0 {
				// 	ctlr.WaitForTxConfirm = true
				// }

			}

			ctrlr, err := transaction.NewController(networkHandler, ks, &account, fromCmdLineFlags)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			if transactionFailure := ctrlr.ExecuteTransaction(toAddress, "", amount, fromShardID, toShardID); transactionFailure != nil {
				fmt.Println(transactionFailure)
				os.Exit(-1)
			}
			fmt.Println(ctrlr.Receipt())
		},
	}

	// TODO Intern do custom flag validation for one address: see https://github.com/spf13/cobra/issues/376
	cmdTransfer.Flags().StringVar(&fromAddress, "from-address", "", "the from address")
	cmdTransfer.Flags().StringVar(&accountName, "account-name", "", "assuming account name to use")
	cmdTransfer.Flags().StringVar(&toAddress, "to-address", "", "the to address")
	cmdTransfer.Flags().Float64Var(&amount, "amount", 0.0, "amount")
	cmdTransfer.Flags().IntVar(&fromShardID, "from-shard", -1, "source shard id")
	cmdTransfer.Flags().IntVar(&toShardID, "to-shard", -1, "target shard id")
	cmdTransfer.PersistentFlags().Uint32Var(&chainID, "chain-id", 0, "What chain ID to signed for, 1 == mainnet, 2 == localnet")
	cmdTransfer.PersistentFlags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "Only waits if non-zero value, in seconds")

	for _, flagName := range [...]string{"from-address", "account-name", "to-address", "amount", "from-shard", "to-shard"} {
		cmdTransfer.MarkFlagRequired(flagName)
	}

	RootCmd.AddCommand(cmdTransfer)
}
