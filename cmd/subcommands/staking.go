package cmd

import (
	"errors"
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
	return nonce.Uint64()
}

func createStakingTransaction(nonce uint64, f staking.StakeMsgFulfiller) (*staking.StakingTransaction, error) {
	gasPrice := big.NewInt(int64(gasPrice))
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(denominations.Nano))

	//TODO: modify the gas limit calculation algorithm
	gasLimit, err := core.IntrinsicGas(nil, false, true)
	if err != nil {
		return nil, err
	}

	stakingTx, err := staking.NewStakingTransaction(nonce, gasLimit, gasPrice, f)
	return stakingTx, err
}

func handleStakingTransaction(stakingTx *staking.StakingTransaction, networkHandler *rpc.HTTPMessenger) error {
	var ks      *keystore.KeyStore
	var acct    *accounts.Account
	var signed  *staking.StakingTransaction
	var err      error

	from := delegatorAddress.String()

	if useLedgerWallet {
		var signerAddr string
		signed, signerAddr, err = ledger.SignStakingTx(stakingTx,  chainName.chainID.Value)
		if err != nil {
			return err
		}

		if strings.Compare(signerAddr, delegatorAddress.String()) != 0 {
			return errors.New("error : delegator address doesn't match with ledger hardware addresss")
		}
	} else {
		ks, acct, err = store.UnlockedKeystore(from, unlockP)
		if err != nil {
			return err
		}
		signed, err = ks.SignStakingTx(*acct, stakingTx, chainName.chainID.Value)
	}

	if err != nil {
		return err
	}

	enc, err := rlp.EncodeToBytes(signed)
	if err != nil {
		return err
	}

	hexSignature := hexutil.Encode(enc)
	reply, err := networkHandler.SendRPC(rpc.Method.SendRawStakingTransaction, []interface{}{hexSignature})
	if err != nil {
		return err
	}
	r, _ := reply["result"].(string)
	fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, r))
	return nil
}

func stakingSubCommands() []*cobra.Command {
	subCmdDelegate := &cobra.Command{
		Use:   "delegate",
		Short: "delegate staking",
		Long: `
Delegating to a validator
`,
		Run: func(cmd *cobra.Command, args []string)  {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				fmt.Println(err)
				return
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveDelegate, staking.Delegate{
					accounts.ParseAddrH(delegatorAddress.String()),
					accounts.ParseAddrH(validatorAddress.String()),
					amt,
				}
			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = handleStakingTransaction(stakingTx, networkHandler)
			if err != nil {
				fmt.Println(err)
				return
			}
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

	for _, flagName := range [...]string{"delegator", "validator", "staking-amount"} {
		subCmdDelegate.MarkFlagRequired(flagName)
	}

	subCmdUnDelegate := &cobra.Command{
		Use:   "undelegate",
		Short: "un-delegate staking",
		Long: `
Remove delegating to a validator
`,
		Run: func(cmd *cobra.Command, args []string)  {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				fmt.Println(err)
				return
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveUndelegate, staking.Undelegate{
					accounts.ParseAddrH(delegatorAddress.String()),
					accounts.ParseAddrH(validatorAddress.String()),
					amt,
				}
			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = handleStakingTransaction(stakingTx, networkHandler)
			if err != nil {
				fmt.Println(err)
				return
			}
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

	for _, flagName := range [...]string{"delegator", "validator", "staking-amount"} {
		subCmdUnDelegate.MarkFlagRequired(flagName)
	}

	subCmdReDelegate := &cobra.Command{
		Use:   "redelegate",
		Short: "re-delegate staking",
		Long: `
Re-delegating to a validator
`,
		Run: func(cmd *cobra.Command, args []string)  {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				fmt.Println(err)
				return
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveUndelegate, staking.Redelegate{
					accounts.ParseAddrH(delegatorAddress.String()),
					accounts.ParseAddrH(validatorSrcAddress.String()),
					accounts.ParseAddrH(validatorDstAddress.String()),
					amt,
				}
			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = handleStakingTransaction(stakingTx, networkHandler)
			if err != nil {
				fmt.Println(err)
				return
			}
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

	for _, flagName := range [...]string{"delegator", "src-validator", "dest-validator", "staking-amount"} {
		subCmdReDelegate.MarkFlagRequired(flagName)
	}

	return []*cobra.Command{subCmdDelegate,
		subCmdUnDelegate,
		subCmdReDelegate,
	}
}

func init() {
	cmdStaking := &cobra.Command{
		Use:   "staking",
		Short: "Delegate, undelegate or redelegate",
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
