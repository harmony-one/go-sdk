package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	common2 "github.com/harmony-one/go-sdk/pkg/common"
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
	amount      string
	fromShardID uint32
	toShardID   uint32
	confirmWait uint32
	chainName   = chainIDWrapper{chainID: &common.Chain.TestNet}
	dryRun      bool
	inputNonce  string
	gasPrice    string
	gasLimit    uint64
	fileFlags   []transferFlags
	txLogs      []transactionLog
	errorFlag   = false
	timeFormat  = "2006-01-02 15:04:05.000000"
)

type transactionLog struct {
	TimeSigned  string
	FromAddress string
	ToAddress   string
	Amount      string
	FromShardID uint32
	ToShardID   uint32
	TxHash      string
	Error       string
}

type transferFlags struct {
	FromAddress *string
	ToAddress   *string
	Amount      *string
	FromShardID *uint32
	ToShardID   *uint32
	UnlockP     *string
	InputNonce  string
	GasPrice    string
	GasLimit    uint64
}

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

func handlerForError(txLog *transactionLog, err error) *transactionLog {
	// Adds error to transactionLog reference
	txLog.Error = time.Now().UTC().Format(timeFormat) + " -- " + err.Error()
	errorFlag = true
	return txLog
}

func stringErrorConstructor(errorString string) transactionLog {
	// Adds error in the form of a string to a new transactionLog and returns it
	errorFlag = true
	return transactionLog{
		Error: time.Now().UTC().Format(timeFormat) + " -- " + errorString,
	}
}

func handlerForTransaction() transactionLog {
	// Executes each individual transaction
	txLog := transactionLog{
		FromAddress: fromAddress.String(),
		ToAddress:   toAddress.String(),
		Amount:      amount,
		FromShardID: fromShardID,
		ToShardID:   toShardID,
	}
	from := fromAddress.String()
	s, err := sharding.Structure(node)
	if err != nil {
		return *handlerForError(&txLog, err)
	}
	err = validation.ValidShardIDs(fromShardID, toShardID, uint32(len(s)))
	if err != nil {
		return *handlerForError(&txLog, err)
	}
	networkHandler, err := handlerForShard(fromShardID, node)
	if err != nil {
		return *handlerForError(&txLog, err)
	}
	var ctrlr *transaction.Controller
	if useLedgerWallet {
		account := accounts.Account{Address: address.Parse(from)}
		ctrlr = transaction.NewController(networkHandler, nil, &account, *chainName.chainID, opts)
	} else {
		ks, acct, err := store.UnlockedKeystore(from, unlockP)
		if err != nil {
			return *handlerForError(&txLog, err)
		}
		ctrlr = transaction.NewController(networkHandler, ks, acct, *chainName.chainID, opts)
	}

	nonce, err := getNonceFromInput(fromAddress.String(), inputNonce, networkHandler)
	if err != nil {
		return *handlerForError(&txLog, err)
	}

	amt, err := common2.NewDecFromString(amount)
	if err != nil {
		return *handlerForError(&txLog, err)
	}

	gPrice, err := common2.NewDecFromString(gasPrice)
	if err != nil {
		return *handlerForError(&txLog, err)
	}

	// Approximate Time of Signature
	txLog.TimeSigned = time.Now().UTC().Format(timeFormat)
	if transactionFailure := ctrlr.ExecuteTransaction(
		toAddress.String(),
		"",
		amt, gPrice,
		nonce, gasLimit,
		int(fromShardID),
		int(toShardID),
	); transactionFailure != nil {
		return *handlerForError(&txLog, transactionFailure)
	}
	txLog.TxHash = *ctrlr.ReceiptHash()
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
	return txLog
}

func handlerForBulkTransactions(index int) transactionLog {
	// Sets flags for a transaction and calls handlerForTransaction()
	// First check that all required flags are present
	if fileFlags[index].FromAddress == nil || fileFlags[index].ToAddress == nil ||
		fileFlags[index].Amount == nil {
		return stringErrorConstructor("FromAddress/ToAddress/Amount are required fields")
	}
	if fileFlags[index].FromShardID == nil || fileFlags[index].ToShardID == nil {
		return stringErrorConstructor("FromShardID/ToShardID are required fields")
	}
	if err := fromAddress.Set(*fileFlags[index].FromAddress); err != nil {
		return stringErrorConstructor(err.Error())
	}
	if err := toAddress.Set(*fileFlags[index].ToAddress); err != nil {
		return stringErrorConstructor(err.Error())
	}
	amount = *fileFlags[index].Amount
	fromShardID = *fileFlags[index].FromShardID
	toShardID = *fileFlags[index].ToShardID
	if fileFlags[index].UnlockP != nil {
		unlockP = *fileFlags[index].UnlockP
	} else {
		unlockP = common.DefaultPassphrase
	}
	inputNonce = fileFlags[index].InputNonce
	if fileFlags[index].GasPrice != "" {
		gasPrice = fileFlags[index].GasPrice
	} else {
		gasPrice = "1"
	}
	if fileFlags[index].GasLimit != 0 {
		gasLimit = fileFlags[index].GasLimit
	} else {
		gasLimit = 21000
	}
	return handlerForTransaction()
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

func getNonceFromInput(addr, inputNonce string, messenger rpc.T) (uint64, error) {
	if inputNonce != "" {
		nonce, err := strconv.ParseUint(inputNonce, 10, 64)
		if err != nil {
			return 0, err
		} else {
			return nonce, nil
		}
	} else {
		return transaction.GetNextNonce(addr, messenger), nil
	}
}

func init() {
	cmdTransfer := &cobra.Command{
		Use:   "transfer",
		Short: "Create and send a transaction",
		Long: `
Create a transaction, sign it, and send off to the Harmony blockchain
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if filepath == "" {
				for _, flagName := range [...]string{"from", "to", "amount", "from-shard", "to-shard"} {
					cmd.MarkFlagRequired(flagName)
				}
			} else {
				data, err := ioutil.ReadFile(filepath)
				if err != nil {
					return err
				}
				err = json.Unmarshal(data, &fileFlags)
				if err != nil {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if filepath == "" {
				txLog := handlerForTransaction()
				txLogs = append(txLogs, txLog)
			} else {
				for i := range fileFlags {
					txLog := handlerForBulkTransactions(i)
					txLogs = append(txLogs, txLog)
				}
			}
			println(common.ToJSONUnsafe(txLogs, true))
			if errorFlag {
				return fmt.Errorf("One or more of your transactions returned an error. Check the log for more information")
			}
			return nil
		},
	}

	cmdTransfer.Flags().Var(&fromAddress, "from", "sender's one address, keystore must exist locally")
	cmdTransfer.Flags().Var(&toAddress, "to", "the destination one address")
	cmdTransfer.Flags().BoolVar(&dryRun, "dry-run", false, "do not send signed transaction")
	cmdTransfer.Flags().StringVar(&amount, "amount", "0", "amount to send (ONE)")
	cmdTransfer.Flags().StringVar(&gasPrice, "gas-price", "1", "gas price to pay (NANO)")
	cmdTransfer.Flags().Uint64Var(&gasLimit, "gas-limit", 21000, "gas limit")
	cmdTransfer.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for tx")
	cmdTransfer.Flags().Uint32Var(&fromShardID, "from-shard", 0, "source shard id")
	cmdTransfer.Flags().Uint32Var(&toShardID, "to-shard", 0, "target shard id")
	cmdTransfer.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	cmdTransfer.Flags().Uint32Var(&confirmWait, "wait-for-confirm", 0, "only waits if non-zero value, in seconds")
	cmdTransfer.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdTransfer.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	RootCmd.AddCommand(cmdTransfer)
}
