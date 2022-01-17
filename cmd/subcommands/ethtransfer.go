package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	rpcEth "github.com/harmony-one/go-sdk/pkg/rpc/eth"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/core"

	"github.com/spf13/cobra"
)

type ethTransferFlags struct {
	FromAddress      *string `json:"from"`
	ToAddress        *string `json:"to"`
	Amount           *string `json:"amount"`
	PassphraseString *string `json:"passphrase-string"`
	PassphraseFile   *string `json:"passphrase-file"`
	InputNonce       *string `json:"nonce"`
	GasPrice         *string `json:"gas-price"`
	GasLimit         *string `json:"gas-limit"`
	StopOnError      bool    `json:"stop-on-error"`
	TrueNonce        bool    `json:"true-nonce"`
}

func ethHandlerForShard(node string) (*rpc.HTTPMessenger, error) {
	return rpc.NewHTTPHandler(node), nil
}

// handlerForTransaction executes a single transaction and fills out the transaction logger accordingly.
//
// Note that the vars need to be set before calling this handler.
func ethHandlerForTransaction(txLog *transactionLog) error {
	from := fromAddress.String()
	networkHandler, err := ethHandlerForShard(node)
	if handlerForError(txLog, err) != nil {
		return err
	}

	var ctrlr *transaction.EthController
	if useLedgerWallet {
		account := accounts.Account{Address: address.Parse(from)}
		ctrlr = transaction.NewEthController(networkHandler, nil, &account, *chainName.chainID, ethOpts)
	} else {
		ks, acct, err := store.UnlockedKeystore(from, passphrase)
		if handlerForError(txLog, err) != nil {
			return err
		}
		ctrlr = transaction.NewEthController(networkHandler, ks, acct, *chainName.chainID, ethOpts)
	}

	nonce, err := getNonce(fromAddress.String(), networkHandler)
	if err != nil {
		return err
	}

	amt, err := common.NewDecFromString(amount)
	if err != nil {
		amtErr := fmt.Errorf("amount %w", err)
		handlerForError(txLog, amtErr)
		return amtErr
	}

	gPrice, err := common.NewDecFromString(gasPrice)
	if err != nil {
		gasErr := fmt.Errorf("gas-price %w", err)
		handlerForError(txLog, gasErr)
		return gasErr
	}

	var gLimit uint64
	if gasLimit == "" {
		gLimit, err = core.IntrinsicGas([]byte(""), false, true, true, false)
		if handlerForError(txLog, err) != nil {
			return err
		}
	} else {
		if strings.HasPrefix(gasLimit, "-") {
			limitErr := fmt.Errorf("gas-limit can not be negative: %s", gasLimit)
			handlerForError(txLog, limitErr)
			return limitErr
		}
		tempLimit, e := strconv.ParseInt(gasLimit, 10, 64)
		if handlerForError(txLog, e) != nil {
			return e
		}
		gLimit = uint64(tempLimit)
	}

	txLog.TimeSigned = time.Now().UTC().Format(timeFormat) // Approximate time of signature
	err = ctrlr.ExecuteEthTransaction(
		nonce, gLimit,
		toAddress.String(),
		amt, gPrice,
		[]byte{},
	)

	if dryRun {
		txLog.RawTxn = ctrlr.RawTransaction()
		txLog.Transaction = make(map[string]interface{})
		_ = json.Unmarshal([]byte(ctrlr.EthTransactionToJSON(false)), &txLog.Transaction)
	} else if txHash := ctrlr.TransactionHash(); txHash != nil {
		txLog.TxHash = *txHash
	}
	txLog.Receipt = ctrlr.Receipt()["result"]
	if err != nil {
		// Report all transaction errors first...
		for _, txError := range ctrlr.TransactionErrors() {
			_ = handlerForError(txLog, txError.Error())
		}
		err = handlerForError(txLog, err)
	}
	if !dryRun && timeout > 0 && txLog.Receipt == nil {
		err = handlerForError(txLog, errors.New("Failed to confirm transaction"))
	}
	return err
}

// ethHandlerForBulkTransactions checks and sets all flags for a transaction
// from the element at index of transferFileFlags, then calls handlerForTransaction.
func ethHandlerForBulkTransactions(txLog *transactionLog, index int) error {
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
		gasPrice = "30000" // Reset to default for subsequent transactions
	}
	if txnFlags.GasLimit != nil {
		gasLimit = *txnFlags.GasLimit
	} else {
		gasLimit = "" // Reset to default for subsequent transactions
	}
	trueNonce = txnFlags.TrueNonce

	return ethHandlerForTransaction(txLog)
}

func ethOpts(ctlr *transaction.EthController) {
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

func init() {
	cmdEthTransfer := &cobra.Command{
		Use:   "eth-transfer",
		Short: "Create and send an Ethereum compatible transaction",
		Args:  cobra.ExactArgs(0),
		Long: `
Create an Ethereum compatible transaction, sign it, and send off to the Harmony blockchain
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if givenFilePath == "" {
				for _, flagName := range [...]string{"from", "to", "amount", "chain-id"} {
					_ = cmd.MarkFlagRequired(flagName)
				}
				if trueNonce && inputNonce != "" {
					return fmt.Errorf("cannot specify nonce when using true on-chain nonce")
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
				for i, batchTx := range transferFileFlags {
					if batchTx.TrueNonce && batchTx.InputNonce != nil {
						return fmt.Errorf("cannot specify nonce when using true on-chain nonce for transaction number %v in batch", i+1)
					}
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc.Method = rpcEth.Method

			if givenFilePath == "" {
				pp, err := getPassphrase()
				if err != nil {
					return err
				}
				passphrase = pp // needed for passphrase assignment used in handler
				txLog := transactionLog{}
				err = ethHandlerForTransaction(&txLog)
				fmt.Println(common.ToJSONUnsafe(txLog, !noPrettyOutput))
				return err
			} else {
				hasError := false
				var txLogs []transactionLog
				for i := range transferFileFlags {
					var txLog transactionLog
					err := ethHandlerForBulkTransactions(&txLog, i)
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

	cmdEthTransfer.Flags().Var(&fromAddress, "from", "sender's one address, keystore must exist locally")
	cmdEthTransfer.Flags().Var(&toAddress, "to", "the destination one address")
	cmdEthTransfer.Flags().BoolVar(&dryRun, "dry-run", false, "do not send signed transaction")
	cmdEthTransfer.Flags().BoolVar(&trueNonce, "true-nonce", false, "send transaction with on-chain nonce")
	cmdEthTransfer.Flags().StringVar(&amount, "amount", "0", "amount to send (ONE)")
	cmdEthTransfer.Flags().StringVar(&gasPrice, "gas-price", "30000", "gas price to pay (NANO)")
	cmdEthTransfer.Flags().StringVar(&gasLimit, "gas-limit", "", "gas limit")
	cmdEthTransfer.Flags().StringVar(&inputNonce, "nonce", "", "set nonce for tx")
	cmdEthTransfer.Flags().StringVar(&targetChain, "chain-id", "", "what chain ID to target")
	cmdEthTransfer.Flags().Uint32Var(&timeout, "timeout", defaultTimeout, "set timeout in seconds. Set to 0 to not wait for confirm")
	cmdEthTransfer.Flags().BoolVar(&userProvidesPassphrase, "passphrase", false, ppPrompt)
	cmdEthTransfer.Flags().StringVar(&passphraseFilePath, "passphrase-file", "", "path to a file containing the passphrase")

	RootCmd.AddCommand(cmdEthTransfer)
}
