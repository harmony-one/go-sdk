package cmd

import (
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	addr oneAddress
	size int64
)

func init() {
	cmdValidator := &cobra.Command{
		Use:   "validator",
		Short: "information about validators",
		Long: `
Look up information about validator information
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdDelegation := &cobra.Command{
		Use:   "delegation",
		Short: "information about delegations",
		Long: `
Look up information about delegation
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdBlockchain := &cobra.Command{
		Use:   "blockchain",
		Short: "Interact with the Harmony.one Blockchain",
		Long: `
Query Harmony's blockchain for completed transaction, historic records
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	accountHistorySubCmd := &cobra.Command{
		Use:   "account-history",
		Short: "Get history of all transactions for given account",
		Long: `
High level information about each transaction for given account
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			type historyParams struct {
				Address   string `json:"address"`
				PageIndex int64  `json:"pageIndex"`
				PageSize  int64  `json:"pageSize"`
				FullTx    bool   `json:"fullTx"`
				TxType    string `json:"txType"`
				Order     string `json:"order"`
			}
			noLatest = true
			params := historyParams{addr.String(), 0, size, true, "", ""}
			return request(rpc.Method.GetTransactionsHistory, []interface{}{params})
		},
	}

	accountHistorySubCmd.Flags().Var(&addr, "address", "address to list transactions for")
	accountHistorySubCmd.Flags().Int64Var(&size, "max-tx", 1000, "max number of transactions to list")

	subCommands := []*cobra.Command{{
		Use:   "block-by-number",
		Short: "Get a harmony blockchain block by block number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetBlockByNumber, []interface{}{args[0], true})
		},
	}, {
		Use:   "known-chains",
		Short: "Print out the known chain-ids",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(common.ToJSONUnsafe(common.AllChainIDs(), !noPrettyOutput))
			return nil
		},
	}, {
		Use:   "protocol-version",
		Short: "The version of the Harmony Protocol",
		Long: `
Query Harmony's blockchain for high level metrics, queries
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(rpc.Method.ProtocolVersion, []interface{}{})
		},
	}, {
		Use:   "transaction-by-hash",
		Short: "Get transaction by hash",
		Args:  cobra.ExactArgs(1),
		Long: `
Inputs of a transaction and r, s, v value of transaction
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetTransactionByHash, []interface{}{args[0]})
		},
	}, {
		Use:   "transaction-receipt",
		Short: "Get information about a finalized transaction",
		Args:  cobra.ExactArgs(1),
		Long: `
High level information about transaction, like blockNumber, blockHash
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetTransactionReceipt, []interface{}{args[0]})
		},
	}, {
		Use:   "median-stake",
		Short: "median stake of top 320 validators with delegations applied stake (pre-epos processing)",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetMedianRawStakeSnapshot, []interface{}{})
		},
	}, {
		Use:   "current-nonce",
		Short: "Current nonce of an account",
		Args:  cobra.ExactArgs(1),
		Long:  `Current nonce number of a one-address`,
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := oneAddress{}
			if err := addr.Set(args[0]); err != nil {
				return err
			}
			return request(rpc.Method.GetTransactionCount, []interface{}{args[0]})
		},
	}, {
		Use:   "pool",
		Short: "Dump a node's transaction pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetPendingTxnsInPool, []interface{}{})
		},
	}, {
		Use:   "latest-header",
		Short: "Get the latest header",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetLatestBlockHeader, []interface{}{})
		},
	}, {
		Use:   "resend-cx",
		Short: "Re-play a cross shard transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.ResendCX, []interface{}{args[0]})
		},
	}, {
		Use:   "utility-metrics",
		Short: "Current utility metrics",
		Long:  `Current staking utility metrics`,
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetCurrentUtilityMetrics, []interface{}{})
		},
	},
		accountHistorySubCmd,
	}

	cmdBlockchain.AddCommand(cmdValidator)
	cmdBlockchain.AddCommand(cmdDelegation)
	cmdValidator.AddCommand(validatorSubCmds[:]...)
	cmdDelegation.AddCommand(delegationSubCmds[:]...)
	cmdBlockchain.AddCommand(subCommands[:]...)
	RootCmd.AddCommand(cmdBlockchain)
}
