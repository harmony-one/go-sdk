package cmd

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/keys"
	"github.com/harmony-one/go-sdk/pkg/ledger"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
	"github.com/harmony-one/harmony/shard"
	staking "github.com/harmony-one/harmony/staking/types"
	"github.com/spf13/cobra"
)

const (
	blsPubKeySize            = 48
	MaxNameLength            = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
	validatorName             string
	validatorIdentity         string
	validatorWebsite          string
	validatorSecurityContact  string
	validatorDetails          string
	commisionRateStr          string
	commisionMaxRateStr       string
	commisionMaxChangeRateStr string
	slotKeyToRemove           string
	slotKeyToAdd              string
	minSelfDelegation         string
	maxTotalDelegation        string
	stakingBlsPubKeys         []string
	delegatorAddress          oneAddress
	validatorAddress          oneAddress
	stakingAmount             string
	oneAsDec                  = numeric.NewDec(denominations.One)
	nanoAsDec                 = numeric.NewDec(denominations.Nano)
)

var (
	errInvalidSelfDelegation           = errors.New("amount value should be between min_self_delegation and max_total_delegation")
	errInvalidTotalDelegation          = errors.New("total delegation can not be bigger than max_total_delegation")
	errMinSelfDelegationTooSmall       = errors.New("min_self_delegation has to be greater than 1 ONE")
	errInvalidMaxTotalDelegation       = errors.New("max_total_delegation can not be less than min_self_delegation")
	errCommissionRateTooLarge          = errors.New("commission rate and change rate can not be larger than max commission rate")
	errInvalidComissionRate            = errors.New("commission rate, change rate and max rate should be within 0-100 percent")
	errInvalidDescFieldName            = errors.New("exceeds maximum length of 70 characters for description field name")
	errInvalidDescFieldIdentity        = errors.New("exceeds maximum length of 3000 characters for description field identity")
	errInvalidDescFieldWebsite         = errors.New("exceeds maximum length of 140 characters for description field website")
	errInvalidDescFieldSecurityContact = errors.New("exceeds maximum length of 140 characters for description field security-contact")
	errInvalidDescFieldDetails         = errors.New("exceeds maximum length of 280 characters for description field details")
)

func createStakingTransaction(nonce uint64, f staking.StakeMsgFulfiller) (*staking.StakingTransaction, error) {
	gPrice := big.NewInt(int64(gasPrice))

	_, payload := f()
	data, err := rlp.EncodeToBytes(payload)
	if err != nil {
		return nil, err
	}

	gasLimit, err := core.IntrinsicGas(data, false, true)
	if err != nil {
		return nil, err
	}

	stakingTx, err := staking.NewStakingTransaction(nonce, gasLimit, gPrice, f)
	return stakingTx, err
}

func handleStakingTransaction(
	stakingTx *staking.StakingTransaction, networkHandler *rpc.HTTPMessenger, signerAddress oneAddress,
) error {
	var (
		ks     *keystore.KeyStore
		acct   *accounts.Account
		signed *staking.StakingTransaction
		err    error
	)

	from := signerAddress.String()

	if useLedgerWallet {
		signerAddr := ""
		signed, signerAddr, err = ledger.SignStakingTx(stakingTx, chainName.chainID.Value)
		if err != nil {
			return err
		}

		if strings.Compare(signerAddr, from) != 0 {
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
	if confirmWait > 0 {
		confirmTx(networkHandler, confirmWait, r)
	} else {
		fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, r))
	}
	return nil
}

func confirmTx(networkHandler *rpc.HTTPMessenger, confirmWaitTime uint32, txHash string) {
	start := int(confirmWaitTime)
	for {
		if start < 0 {
			fmt.Println("Could not confirm", txHash, "even after", confirmWaitTime, "seconds")
			fmt.Println("Try increasing the wait-for-confirm")
			return
		}
		r, _ := networkHandler.SendRPC(rpc.Method.GetTransactionReceipt, []interface{}{txHash})
		if r["result"] != nil {
			fmt.Println(common.ToJSONUnsafe(r, true))
			return
		}
		time.Sleep(time.Second * 2)
		start = start - 2
	}
}

func delegationAmountSanityCheck(minSelfDelegation numeric.Dec, maxTotalDelegation numeric.Dec, amount *numeric.Dec) error {
	// MinSelfDelegation must be >= 1 ONE
	if minSelfDelegation.LT(oneAsDec) {
		return errMinSelfDelegationTooSmall
	}

	// MaxTotalDelegation must not be less than MinSelfDelegation
	if !maxTotalDelegation.Equal(numeric.ZeroDec()) &&
		maxTotalDelegation.LT(minSelfDelegation) {
		return errInvalidMaxTotalDelegation
	}

	// Amount must be >= MinSelfDelegation
	if amount != nil && amount.LT(minSelfDelegation) &&
		(maxTotalDelegation.Equal(numeric.ZeroDec()) ||
			(!maxTotalDelegation.Equal(numeric.ZeroDec()) && amount.GT(maxTotalDelegation))) {
		return errInvalidSelfDelegation
	}

	return nil
}

func rateSanityCheck(rate numeric.Dec, maxRate numeric.Dec, maxChangeRate numeric.Dec) error {
	hundredPercent := numeric.NewDec(1)
	zeroPercent := numeric.NewDec(0)

	if rate.LT(zeroPercent) || rate.GT(hundredPercent) {
		return errInvalidComissionRate
	}

	if maxRate.LT(zeroPercent) || maxRate.GT(hundredPercent) {
		return errInvalidComissionRate
	}

	if maxChangeRate.LT(zeroPercent) || maxChangeRate.GT(hundredPercent) {
		return errInvalidComissionRate
	}

	if rate.GT(maxRate) {
		return errCommissionRateTooLarge
	}

	if maxChangeRate.GT(maxRate) {
		return errCommissionRateTooLarge
	}

	return nil
}

func ensureLength(d staking.Description) (staking.Description, error) {
	if len(d.Name) > MaxNameLength {
		return d, errInvalidDescFieldName
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, errInvalidDescFieldIdentity
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, errInvalidDescFieldWebsite
	}
	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, errInvalidDescFieldSecurityContact
	}
	if len(d.Details) > MaxDetailsLength {
		return d, errInvalidDescFieldDetails
	}

	return d, nil
}

func stakingSubCommands() []*cobra.Command {

	subCmdNewValidator := &cobra.Command{
		Use:   "create-validator",
		Short: "create a new validator",
		Long: `
Create a new validator"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			commisionRate, err := numeric.NewDecFromStr(commisionRateStr)
			if err != nil {
				return err
			}

			commisionMaxRate, err := numeric.NewDecFromStr(commisionMaxRateStr)
			if err != nil {
				return err
			}

			commisionMaxChangeRate, err := numeric.NewDecFromStr(commisionMaxChangeRateStr)
			if err != nil {
				return err
			}

			blsPubKeys := make([]shard.BlsPublicKey, len(stakingBlsPubKeys))
			for i := 0; i < len(stakingBlsPubKeys); i++ {
				blsPubKey := new(bls.PublicKey)
				err = blsPubKey.DeserializeHexStr(strings.TrimPrefix(stakingBlsPubKeys[i], "0x"))
				if err != nil {
					return err
				}

				blsPubKeys[i].FromLibBLSPublicKey(blsPubKey)
			}

			blsSigs, err := keys.VerifyBLSKeys(stakingBlsPubKeys)
			if err != nil {
				return err
			}

			amt, e0 := numeric.NewDecFromStr(stakingAmount)
			minSelfDel, e1 := numeric.NewDecFromStr(minSelfDelegation)
			maxTotalDel, e2 := numeric.NewDecFromStr(maxTotalDelegation)

			if e0 != nil {
				return e0
			}
			if e1 != nil {
				return e1
			}
			if e2 != nil {
				return e2
			}

			amt = amt.Mul(oneAsDec)
			minSelfDel = minSelfDel.Mul(oneAsDec)
			maxTotalDel = maxTotalDel.Mul(oneAsDec)

			err = delegationAmountSanityCheck(minSelfDel, maxTotalDel, &amt)
			if err != nil {
				return err
			}

			err = rateSanityCheck(commisionRate, commisionMaxRate, commisionMaxChangeRate)
			if err != nil {
				return err
			}

			desc, err := ensureLength(staking.Description{
				validatorName,
				validatorIdentity,
				validatorWebsite,
				validatorSecurityContact,
				validatorDetails,
			})
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				return staking.DirectiveCreateValidator, staking.CreateValidator{
					address.Parse(validatorAddress.String()),
					&desc,
					staking.CommissionRates{
						commisionRate,
						commisionMaxRate,
						commisionMaxChangeRate},
					minSelfDel.RoundInt(),
					maxTotalDel.RoundInt(),
					blsPubKeys,
					blsSigs,
					amt.RoundInt(),
				}
			}

			nonce, err := getNonceFromInput(validatorAddress.String(), inputNonce, networkHandler)
			if err != nil {
				return err
			}
			stakingTx, err := createStakingTransaction(nonce, delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, validatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdNewValidator.Flags().StringVar(&validatorName, "name", "", "validator's name")
	subCmdNewValidator.Flags().StringVar(&validatorIdentity, "identity", "", "validator's identity")
	subCmdNewValidator.Flags().StringVar(&validatorWebsite, "website", "", "validator's website")
	subCmdNewValidator.Flags().StringVar(&validatorSecurityContact, "security-contact", "", "validator's security contact")
	subCmdNewValidator.Flags().StringVar(&validatorDetails, "details", "", "validator's details")
	subCmdNewValidator.Flags().StringVar(&commisionRateStr, "rate", "", "commission rate")
	subCmdNewValidator.Flags().StringVar(&commisionMaxRateStr, "max-rate", "", "commision max rate")
	subCmdNewValidator.Flags().StringVar(&commisionMaxChangeRateStr, "max-change-rate", "", "commission max change amount")
	subCmdNewValidator.Flags().StringVar(&minSelfDelegation, "min-self-delegation", "0.0", "minimal self delegation")
	subCmdNewValidator.Flags().StringVar(&maxTotalDelegation, "max-total-delegation", "0.0", "maximal total delegation")
	subCmdNewValidator.Flags().Var(&validatorAddress, "validator-addr", "validator's staking address")
	subCmdNewValidator.Flags().StringSliceVar(
		&stakingBlsPubKeys, "bls-pubkeys", []string{}, "validator's list of public BLS key addresses",
	)
	subCmdNewValidator.Flags().StringVar(&stakingAmount, "amount", "0.0", "staking amount")
	subCmdNewValidator.Flags().Uint64Var(&gasPrice, "gas-price", 1, "gas price to pay")
	subCmdNewValidator.Flags().IntVar(&gasLimit, "gas-limit", 21000, "gas limit")
	subCmdNewValidator.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for transaction")
	subCmdNewValidator.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdNewValidator.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	subCmdNewValidator.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"name", "identity", "website", "security-contact", "details", "rate", "max-rate",
		"max-change-rate", "min-self-delegation", "max-total-delegation", "validator-addr", "bls-pubkeys", "amount"} {
		subCmdNewValidator.MarkFlagRequired(flagName)
	}

	subCmdEditValidator := &cobra.Command{
		Use:   "edit-validator",
		Short: "edit a validator",
		Long:  "Edit an existing validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			commisionRate, err := numeric.NewDecFromStr(commisionRateStr)
			if err != nil {
				return err
			}

			blsPubKeyRemove := new(bls.PublicKey)
			err = blsPubKeyRemove.DeserializeHexStr(strings.TrimPrefix(slotKeyToRemove, "0x"))
			if err != nil {
				return err
			}

			shardPubKeyRemove := shard.BlsPublicKey{}
			shardPubKeyRemove.FromLibBLSPublicKey(blsPubKeyRemove)

			blsPubKeyAdd := new(bls.PublicKey)
			err = blsPubKeyAdd.DeserializeHexStr(strings.TrimPrefix(slotKeyToAdd, "0x"))
			if err != nil {
				return err
			}

			shardPubKeyAdd := shard.BlsPublicKey{}
			shardPubKeyAdd.FromLibBLSPublicKey(blsPubKeyAdd)

			sigBls, err := keys.VerifyBLS(strings.TrimPrefix(slotKeyToAdd, "0x"))
			if err != nil {
				return err
			}

			minSelfDel, e1 := numeric.NewDecFromStr(minSelfDelegation)
			maxTotalDel, e2 := numeric.NewDecFromStr(maxTotalDelegation)

			if e1 != nil {
				return e1
			}
			if e2 != nil {
				return e2
			}

			minSelfDel = minSelfDel.Mul(oneAsDec)
			maxTotalDel = maxTotalDel.Mul(oneAsDec)

			err = delegationAmountSanityCheck(minSelfDel, maxTotalDel, nil)
			if err != nil {
				return err
			}

			desc, err := ensureLength(staking.Description{
				validatorName,
				validatorIdentity,
				validatorWebsite,
				validatorSecurityContact,
				validatorDetails,
			})
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				return staking.DirectiveEditValidator, staking.EditValidator{
					address.Parse(validatorAddress.String()),
					&desc,
					&commisionRate,
					minSelfDel.RoundInt(),
					maxTotalDel.RoundInt(),
					&shardPubKeyRemove,
					&shardPubKeyAdd,
					sigBls,
				}

			}

			nonce, err := getNonceFromInput(validatorAddress.String(), inputNonce, networkHandler)
			if err != nil {
				return err
			}
			stakingTx, err := createStakingTransaction(nonce, delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, validatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdEditValidator.Flags().StringVar(&validatorName, "name", "", "validator's name")
	subCmdEditValidator.Flags().StringVar(&validatorIdentity, "identity", "", "validator's identity")
	subCmdEditValidator.Flags().StringVar(&validatorWebsite, "website", "", "validator's website")
	subCmdEditValidator.Flags().StringVar(&validatorSecurityContact, "security-contact", "", "validator's security contact")
	subCmdEditValidator.Flags().StringVar(&validatorDetails, "details", "", "validator's details")
	subCmdEditValidator.Flags().StringVar(&commisionRateStr, "rate", "", "commission rate")
	subCmdEditValidator.Flags().StringVar(&minSelfDelegation, "min-self-delegation", "0.0", "minimal self delegation")
	subCmdEditValidator.Flags().StringVar(&maxTotalDelegation, "max-total-delegation", "0.0", "maximal total delegation")
	subCmdEditValidator.Flags().Var(&validatorAddress, "validator-addr", "validator's staking address")
	subCmdEditValidator.Flags().StringVar(&slotKeyToAdd, "add-bls-key", "", "add BLS pubkey to slot")
	subCmdEditValidator.Flags().StringVar(&slotKeyToRemove, "remove-bls-key", "", "remove BLS pubkey from slot")

	subCmdEditValidator.Flags().Uint64Var(&gasPrice, "gas-price", 1, "gas price to pay")
	subCmdEditValidator.Flags().IntVar(&gasLimit, "gas-limit", 21000, "gas limit")
	subCmdEditValidator.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for transaction")
	subCmdEditValidator.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdEditValidator.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	subCmdEditValidator.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{
		"name", "identity", "website", "security-contact", "details", "rate",
		"min-self-delegation", "max-total-delegation", "validator-addr",
		"remove-bls-key", "add-bls-key",
	} {
		subCmdEditValidator.MarkFlagRequired(flagName)
	}

	subCmdDelegate := &cobra.Command{
		Use:   "delegate",
		Short: "delegating to a validator",
		Long: `
Delegating to a validator
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amt, _ := numeric.NewDecFromStr(stakingAmount)
				amt = amt.Mul(oneAsDec)
				return staking.DirectiveDelegate, staking.Delegate{
					address.Parse(delegatorAddress.String()),
					address.Parse(validatorAddress.String()),
					amt.RoundInt(),
				}
			}

			nonce, err := getNonceFromInput(delegatorAddress.String(), inputNonce, networkHandler)
			if err != nil {
				return err
			}
			stakingTx, err := createStakingTransaction(nonce, delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, delegatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdDelegate.Flags().Var(&delegatorAddress, "delegator-addr", "delegator's address")
	subCmdDelegate.Flags().Var(&validatorAddress, "validator-addr", "validator's address")
	subCmdDelegate.Flags().StringVar(&stakingAmount, "amount", "0.0", "staking amount")
	subCmdDelegate.Flags().Uint64Var(&gasPrice, "gas-price", 1, "gas price to pay")
	subCmdDelegate.Flags().IntVar(&gasLimit, "gas-limit", 21000, "gas limit")
	subCmdDelegate.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for transaction")
	subCmdDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdDelegate.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	subCmdDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"delegator-addr", "validator-addr", "amount"} {
		subCmdDelegate.MarkFlagRequired(flagName)
	}

	subCmdUnDelegate := &cobra.Command{
		Use:   "undelegate",
		Short: "removing delegation responsibility",
		Long: `
 Removing delegation responsibility
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amt, _ := numeric.NewDecFromStr(stakingAmount)
				amt = amt.Mul(oneAsDec)

				return staking.DirectiveUndelegate, staking.Undelegate{
					address.Parse(delegatorAddress.String()),
					address.Parse(validatorAddress.String()),
					amt.RoundInt(),
				}
			}

			nonce, err := getNonceFromInput(delegatorAddress.String(), inputNonce, networkHandler)
			if err != nil {
				return err
			}
			stakingTx, err := createStakingTransaction(nonce, delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, delegatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdUnDelegate.Flags().Var(&delegatorAddress, "delegator-addr", "delegator's address")
	subCmdUnDelegate.Flags().Var(&validatorAddress, "validator-addr", "source validator's address")
	subCmdUnDelegate.Flags().StringVar(&stakingAmount, "amount", "0.0", "staking amount")
	subCmdUnDelegate.Flags().Uint64Var(&gasPrice, "gas-price", 1, "gas price to pay")
	subCmdUnDelegate.Flags().IntVar(&gasLimit, "gas-limit", 21000, "gas limit")
	subCmdUnDelegate.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for transaction")
	subCmdUnDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdUnDelegate.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	subCmdUnDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"delegator-addr", "validator-addr", "amount"} {
		subCmdUnDelegate.MarkFlagRequired(flagName)
	}

	subCmdCollectRewards := &cobra.Command{
		Use:   "collect-rewards",
		Short: "collect token rewards",
		Long: `
Collect token rewards
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				return staking.DirectiveCollectRewards, staking.CollectRewards{
					address.Parse(delegatorAddress.String()),
				}
			}

			nonce, err := getNonceFromInput(delegatorAddress.String(), inputNonce, networkHandler)
			if err != nil {
				return err
			}
			stakingTx, err := createStakingTransaction(nonce, delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, delegatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdCollectRewards.Flags().Var(&delegatorAddress, "delegator-addr", "delegator's address")
	subCmdCollectRewards.Flags().Uint64Var(&gasPrice, "gas-price", 1, "gas price to pay")
	subCmdCollectRewards.Flags().IntVar(&gasLimit, "gas-limit", 21000, "gas limit")
	subCmdCollectRewards.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for tx")
	subCmdCollectRewards.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdCollectRewards.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	subCmdCollectRewards.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"delegator-addr"} {
		subCmdCollectRewards.MarkFlagRequired(flagName)
	}

	return []*cobra.Command{
		subCmdNewValidator,
		subCmdEditValidator,
		subCmdDelegate,
		subCmdUnDelegate,
		subCmdCollectRewards,
	}
}

func init() {
	cmdStaking := &cobra.Command{
		Use:   "staking",
		Short: "newvalidator, editvalidator, delegate, undelegate or redelegate",
		Long: `
Create a staking transaction, sign it, and send off to the Harmony blockchain
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdStaking.AddCommand(stakingSubCommands()...)
	RootCmd.AddCommand(cmdStaking)
}
