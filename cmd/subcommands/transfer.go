package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	c "github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/go-sdk/pkg/validation"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/core"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const defaultTimeout = 40

var (
	fromAddress       oneAddress
	toAddress         oneAddress
	amount            string
	fromShardID       uint32
	toShardID         uint32
	chainName         = chainIDWrapper{chainID: &common.Chain.TestNet}
	dryRun            bool
	inputNonce        string
	gasPrice          string
	gasLimit          string
	transferFileFlags []transferFlags
	timeout           uint32
	timeFormat        = "2006-01-02 15:04:05.000000"
)

type transactionLog struct {
	TxHash      string      `json:"transaction-hash,omitempty"`
	Transaction interface{} `json:"transaction,omitempty"`
	Receipt     interface{} `json:"blockchain-receipt,omitempty"`
	RawTxn      string      `json:"raw-transaction,omitempty"`
	Errors      []string    `json:"errors,omitempty"`
	TimeSigned  string      `json:"time-signed-utc,omitempty"`
}

type transferFlags struct {
	FromAddress      *string `json:"from"`
	ToAddress        *string `json:"to"`
	Amount           *string `json:"amount"`
	FromShardID      *string `json:"from-shard"`
	ToShardID        *string `json:"to-shard"`
	PassphraseString *string `json:"passphrase-string"`
	PassphraseFile   *string `json:"passphrase-file"`
	InputNonce       *string `json:"nonce"`
	GasPrice         *string `json:"gas-price"`
	GasLimit         *string `json:"gas-limit"`
	StopOnError      bool    `json:"stop-on-error"`
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

// handlerForError sets the error in the transaction logger to the given error.
// It returns the given error for convenience.
func handlerForError(txLog *transactionLog, err error) error {
	if err != nil {
		errorString := fmt.Sprintf("[%s] %s", time.Now().UTC().Format(timeFormat), err.Error())
		txLog.Errors = append(txLog.Errors, errorString)
	}
	return err
}

// handlerForTransaction executes a single transaction and fills out the transaction logger accordingly.
//
// Note that the vars need to be set before calling this handler.
func handlerForTransaction(txLog *transactionLog) error {
	from := fromAddress.String()
	s, err := sharding.Structure(node)
	if handlerForError(txLog, err) != nil {
		return err
	}
	err = validation.ValidShardIDs(fromShardID, toShardID, uint32(len(s)))
	if handlerForError(txLog, err) != nil {
		return err
	}
	networkHandler, err := handlerForShard(fromShardID, node)
	if handlerForError(txLog, err) != nil {
		return err
	}

	var ctrlr *transaction.Controller
	if useLedgerWallet {
		account := accounts.Account{Address: address.Parse(from)}
		ctrlr = transaction.NewController(networkHandler, nil, &account, *chainName.chainID, opts)
	} else {
		ks, acct, err := store.UnlockedKeystore(from, passphrase)
		if handlerForError(txLog, err) != nil {
			return err
		}
		ctrlr = transaction.NewController(networkHandler, ks, acct, *chainName.chainID, opts)
	}

	nonce, err := getNonceFromInput(fromAddress.String(), inputNonce, networkHandler)
	if handlerForError(txLog, err) != nil {
		return err
	}
	amt, err := c.NewDecFromString(amount)
	if handlerForError(txLog, err) != nil {
		return err
	}
	gPrice, err := c.NewDecFromString(gasPrice)
	if handlerForError(txLog, err) != nil {
		return err
	}

	var gLimit uint64
	if gasLimit == "" {
		gLimit, err = core.IntrinsicGas([]byte(""), false, true, false)
		if handlerForError(txLog, err) != nil {
			return err
		}
	} else {
		tempLimit, e := strconv.ParseInt(gasLimit, 10, 64)
		if handlerForError(txLog, e) != nil {
			return e
		}
		gLimit = uint64(tempLimit)
	}

	txLog.TimeSigned = time.Now().UTC().Format(timeFormat) // Approximate time of signature
	err = ctrlr.ExecuteTransaction(
		nonce, gLimit,
		toAddress.String(),
		int(fromShardID), int(toShardID),
		amt, gPrice,
		[]byte{},
	)

	if dryRun {
		txLog.RawTxn = ctrlr.RawTransaction()
		txLog.Transaction = make(map[string]interface{})
		_ = json.Unmarshal([]byte(ctrlr.TransactionToJSON(false)), &txLog.Transaction)
	} else if txHash := ctrlr.TransactionHash(); txHash != nil {
		txLog.TxHash = *txHash
	}
	txLog.Receipt = ctrlr.Receipt()["result"]
	if timeout > 0 && txLog.Receipt == nil {
		// Report all transaction errors first...
		for _, txError := range ctrlr.TransactionErrors() {
			_ = handlerForError(txLog, txError.Error())
		}
		if err == nil {
			err = handlerForError(txLog, err)
		} else {
			err = handlerForError(txLog, errors.New("Failed to confirm transaction"))
		}
	}
	return err
}

// handlerForBulkTransactions checks and sets all flags for a transaction
// from the element at index of transferFileFlags, then calls handlerForTransaction.
func handlerForBulkTransactions(txLog *transactionLog, index int) error {
	txnFlags := transferFileFlags[index]

	// Check for required fields.
	if txnFlags.FromAddress == nil || txnFlags.ToAddress == nil || txnFlags.Amount == nil {
		return handlerForError(txLog, errors.New("FromAddress/ToAddress/Amount are required fields"))
	}
	if txnFlags.FromShardID == nil || txnFlags.ToShardID == nil {
		return handlerForError(txLog, errors.New("FromShardID/ToShardID are required fields"))
	}

	// Set required fields.
	err := fromAddress.Set(*txnFlags.FromAddress)
	if handlerForError(txLog, err) != nil {
		return err
	}
	err = toAddress.Set(*txnFlags.ToAddress)
	if handlerForError(txLog, err) != nil {
		return err
	}
	amount = *txnFlags.Amount
	fromShard, err := strconv.ParseUint(*txnFlags.FromShardID, 10, 64)
	if handlerForError(txLog, err) != nil {
		return err
	}
	fromShardID = uint32(fromShard)
	toShard, err := strconv.ParseUint(*txnFlags.ToShardID, 10, 64)
	if handlerForError(txLog, err) != nil {
		return err
	}
	toShardID = uint32(toShard)

	// Set optional fields.
	if txnFlags.PassphraseFile != nil {
		passphraseFilePath = *txnFlags.PassphraseFile
		passphrase, err = getPassphrase()
		if handlerForError(txLog, err) != nil {
			return err
		}
	} else if txnFlags.PassphraseString != nil {
		passphrase = *txnFlags.PassphraseString
	} else {
		passphrase = common.DefaultPassphrase
	}
	if txnFlags.InputNonce != nil {
		inputNonce = *txnFlags.InputNonce
	} else {
		inputNonce = "" // Reset to default for subsequent transactions
	}
	if txnFlags.GasPrice != nil {
		gasPrice = *txnFlags.GasPrice
	} else {
		gasPrice = "1" // Reset to default for subsequent transactions
	}
	if txnFlags.GasLimit != nil {
		gasLimit = *txnFlags.GasLimit
	} else {
		gasLimit = "" // Reset to default for subsequent transactions
	}

	return handlerForTransaction(txLog)
}

func opts(ctlr *transaction.Controller) {
	if dryRun {
		ctlr.Behavior.DryRun = true
	}
	if useLedgerWallet {
		ctlr.Behavior.SigningImpl = transaction.Ledger
	}
	if timeout > 0 {
		ctlr.Behavior.ConfirmationWaitTime = timeout
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

type Error struct {
	Msg    string `json:"error-message"`
	TxHash string `json:"tx-hash-id"`
}

func reportError(method string, txHash string) error {
	success, failure := rpc.Request(method, node, []interface{}{})
	if failure != nil {
		return failure
	}
	asJSON, _ := json.Marshal(success["result"])
	var errs []Error
	json.Unmarshal(asJSON, &errs)
	for _, err := range errs {
		if err.TxHash == txHash {
			fmt.Println(fmt.Errorf("error: %s", err.Msg))
			return nil
		}
	}
	return errors.New("could not find error msg")
}

func init() {
	cmdTransfer := &cobra.Command{
		Use:   "transfer",
		Short: "Create and send a transaction",
		Long: `
Create a transaction, sign it, and send off to the Harmony blockchain
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if givenFilePath == "" {
				for _, flagName := range [...]string{"from", "to", "amount", "from-shard", "to-shard"} {
					_ = cmd.MarkFlagRequired(flagName)
				}
			} else {
				data, err := ioutil.ReadFile(givenFilePath)
				if err != nil {
					return err
				}
				err = json.Unmarshal(data, &transferFileFlags)
				if err != nil {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if givenFilePath == "" {
				pp, err := getPassphrase()
				if err != nil {
					return err
				}
				passphrase = pp // needed for passphrase assignment used in handler
				txLog := transactionLog{}
				err = handlerForTransaction(&txLog)
				fmt.Println(common.ToJSONUnsafe(txLog, !noPrettyOutput))
				return err
			} else {
				hasError := false
				var txLogs []transactionLog
				for i := range transferFileFlags {
					var txLog transactionLog
					err := handlerForBulkTransactions(&txLog, i)
					txLogs = append(txLogs, txLog)
					if err != nil {
						hasError = true
						if transferFileFlags[i].StopOnError {
							break
						}
					}
				}
				fmt.Println(common.ToJSONUnsafe(txLogs, true))
				if hasError {
					return fmt.Errorf("one or more of your transactions returned an error " +
						"-- check the log for more information")
				} else {
					return nil
				}
			}
		},
	}

	cmdTransfer.Flags().Var(&fromAddress, "from", "sender's one address, keystore must exist locally")
	cmdTransfer.Flags().Var(&toAddress, "to", "the destination one address")
	cmdTransfer.Flags().BoolVar(&dryRun, "dry-run", false, "do not send signed transaction")
	cmdTransfer.Flags().StringVar(&amount, "amount", "0", "amount to send (ONE)")
	cmdTransfer.Flags().StringVar(&gasPrice, "gas-price", "1", "gas price to pay (NANO)")
	cmdTransfer.Flags().StringVar(&gasLimit, "gas-limit", "", "gas limit")
	cmdTransfer.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for tx")
	cmdTransfer.Flags().Uint32Var(&fromShardID, "from-shard", 0, "source shard id")
	cmdTransfer.Flags().Uint32Var(&toShardID, "to-shard", 0, "target shard id")
	cmdTransfer.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	cmdTransfer.Flags().Uint32Var(&timeout, "timeout", defaultTimeout, "set timeout in seconds. Set to 0 to not wait for confirm")
	cmdTransfer.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdTransfer.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	RootCmd.AddCommand(cmdTransfer)
}
