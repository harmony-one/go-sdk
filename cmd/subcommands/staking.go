package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/ledger"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
	"strings"

	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/common/denominations"
	staking "github.com/harmony-one/harmony/staking/types"
	"github.com/spf13/cobra"
	"math/big"
)

var (
	delegatorAddress    oneAddress
	validatorAddress    oneAddress
	validatorSrcAddress oneAddress
	validatorDstAddress oneAddress
	senderAddress       oneAddress
	stakingAmount       float64
)

func getNextNonce(messenger rpc.T) uint64 {
	transactionCountRPCReply, err :=
		messenger.SendRPC(rpc.Method.GetTransactionCount, []interface{}{accounts.ParseAddrH(delegatorAddress.String()), "latest"})

	if err != nil {
		return 0
	}

	transactionCount, _ := transactionCountRPCReply["result"].(string)
	nonce, _ := big.NewInt(0).SetString(transactionCount[2:], 16)
	fmt.Println("Nonce = ", nonce.String())
	return nonce.Uint64()
}

func createDelegateStakingTransaction(nonce uint64) (*staking.StakingTransaction, error) {
	amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))
	gasPrice := big.NewInt(int64(gasPrice))
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(denominations.Nano))

	//TODO: modify the gas limit calculation algorithm
	gasLimit, err := core.IntrinsicGas(nil, false, true)
	if err != nil {
		return nil, err
	}

	stakePayloadMaker := func() (staking.Directive, interface{}) {
		return staking.DirectiveDelegate, staking.Delegate{
			accounts.ParseAddrH(delegatorAddress.String()),
			accounts.ParseAddrH(validatorAddress.String()),
			amt,
		}
	}

	stakingTx, err := staking.NewStakingTransaction(nonce, gasLimit, gasPrice, stakePayloadMaker)
	return stakingTx, err
}

func stakingSubCommands() []*cobra.Command {
	subCmdDelegate := &cobra.Command{
		Use:   "delegate",
		Short: "delegate staking",
		Long: `
Delegating to a validator
`,
		Run: func(cmd *cobra.Command, args []string)  {
			var ks      *keystore.KeyStore
			var acct    *accounts.Account
			var signed  *staking.StakingTransaction

			networkHandler, err := handlerForShard(0, node)

			from := delegatorAddress.String()
			stakingTx, err := createDelegateStakingTransaction(getNextNonce(networkHandler))
			if err != nil {
				fmt.Println(err)
				return
			}

			if useLedgerWallet {
				var signerAddr string
				signed, signerAddr, err = ledger.SignStakingTx(stakingTx,  chainName.chainID.Value)
				if err != nil {
					fmt.Println(err)
					return
				}

				if strings.Compare(signerAddr, delegatorAddress.String()) != 0 {
					fmt.Println("signature verification failed : delegator address doesn't match with ledger hardware addresss")
					return
				}
			} else {
				ks, acct, err = store.UnlockedKeystore(from, unlockP)
				if err != nil {
					fmt.Println(err)
					return
				}
				signed, err = ks.SignStakingTx(*acct, stakingTx, chainName.chainID.Value)
			}

			if err != nil {
				fmt.Println(err)
				return
			}

			enc, err := rlp.EncodeToBytes(signed)
			if err != nil {
				fmt.Println(err)
				return
			}

			hexSignature := hexutil.Encode(enc)
			reply, err := networkHandler.SendRPC(rpc.Method.SendRawStakingTransaction, []interface{}{hexSignature})
			if err != nil {
				fmt.Println(err)
				return
			}
			r, _ := reply["result"].(string)
			fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, r))
		},
	}

	subCmdDelegate.Flags().Var(&delegatorAddress, "delegator", "delegator's address")
	subCmdDelegate.Flags().Var(&validatorAddress, "validator", "validator's address")
	subCmdDelegate.Flags().Float64Var(&stakingAmount, "staking-amount", 0.0, "staking amount")
	subCmdDelegate.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	subCmdUnDelegate := &cobra.Command{
		Use:   "undelegate",
		Short: "un-delegate staking",
		Long: `
Remove delegating to a validator
`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string)  {
		},
	}

	subCmdUnDelegate.Flags().Var(&delegatorAddress, "delegator", "delegator's address")
	subCmdUnDelegate.Flags().Var(&validatorAddress, "validator", "source validator's address")
	subCmdUnDelegate.Flags().Float64Var(&stakingAmount, "staking-amount", 0.0, "staking amount")
	subCmdUnDelegate.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdUnDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdUnDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)
	subCmdReDelegate := &cobra.Command{
		Use:   "redelegate",
		Short: "re-delegate staking",
		Long: `
Re-delegating to a validator
`,
		Run: func(cmd *cobra.Command, args []string)  {
		},
	}

	subCmdReDelegate.Flags().Var(&delegatorAddress, "delegator", "delegator's address")
	subCmdReDelegate.Flags().Float64Var(&stakingAmount, "staking-amount", 0.0, "staking amount")
	subCmdReDelegate.Flags().Var(&validatorSrcAddress, "src-validator", "source validator's address")
	subCmdReDelegate.Flags().Var(&validatorDstAddress, "dest-validator", "destination validator's address")
	subCmdReDelegate.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdReDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdReDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	return []*cobra.Command{subCmdDelegate,
		subCmdUnDelegate,
		subCmdReDelegate,
	}
}

func init() {
	cmdStaking := &cobra.Command{
		Use:   "staking",
		Short: "Stake, delegate or undelegate",
		Long: `
Create a staking transaction, sign it, and send off to the Harmony blockchain
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmdStaking.AddCommand(stakingSubCommands()...)
	RootCmd.AddCommand(cmdStaking)
}
