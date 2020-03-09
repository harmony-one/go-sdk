package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	color "github.com/fatih/color"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	verbose         bool
	useLedgerWallet bool
	noLatest        bool
	noPrettyOutput  bool
	node            string
	keyStoreDir     string
	givenFilePath   string
	request         = func(method string, params []interface{}) error {
		if !noLatest {
			params = append(params, "latest")
		}
		success, failure := rpc.Request(method, node, params)
		if failure != nil {
			return failure
		}
		asJSON, _ := json.Marshal(success)
		if noPrettyOutput {
			fmt.Println(string(asJSON))
			return nil
		}
		fmt.Println(common.JSONPrettyFormat(string(asJSON)))
		return nil
	}
	// RootCmd is single entry point of the CLI
	RootCmd = &cobra.Command{
		Use:          "hmy",
		Short:        "Harmony blockchain",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				common.EnableAllVerbose()
			}
			if strings.HasPrefix(node, "https://") || strings.HasPrefix(node, "http://") ||
				strings.HasPrefix(node, "ws://") {
				//No op, already has protocol, respect protocol default ports.
			} else if strings.HasPrefix(node, "api") || strings.HasPrefix(node, "ws") {
				node = "https://" + node
			} else {
				switch URLcomponents := strings.Split(node, ":"); len(URLcomponents) {
				case 1:
					node = "http://" + node + ":9500"
				case 2:
					node = "http://" + node
				default:
					node = node
				}
			}

			if targetChain == "" {
				if strings.Contains(node, ".t.") {
					chainName = chainIDWrapper{chainID: &common.Chain.MainNet}
				} else if strings.Contains(node, ".b.") {
					chainName = chainIDWrapper{chainID: &common.Chain.TestNet}
				} else if strings.Contains(node, ".os.") {
					chainName = chainIDWrapper{chainID: &common.Chain.PangaeaNet}
				} else if strings.Contains(node, ".ps.") {
					chainName = chainIDWrapper{chainID: &common.Chain.PartnerNet}
				} else if strings.Contains(node, ".stn.") {
					chainName = chainIDWrapper{chainID: &common.Chain.StressNet}
				} else {
					chainName = chainIDWrapper{chainID: &common.Chain.TestNet}
				}
			} else {
				chain, err := common.StringToChainID(targetChain)
				if err != nil {
					return err
				}
				chainName = chainIDWrapper{chainID: chain}
			}

			return nil
		},
		Long: fmt.Sprintf(`
CLI interface to the Harmony blockchain

%s`, g("Invoke 'hmy cookbook' for examples of the most common, important usages")),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
)

func init() {
	vS := "dump out debug information, same as env var HMY_ALL_DEBUG=true"
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, vS)
	RootCmd.PersistentFlags().StringVarP(&node, "node", "n", defaultNodeAddr, "<host>")
	RootCmd.PersistentFlags().BoolVar(
		&noLatest, "no-latest", false, "Do not add 'latest' to RPC params",
	)
	RootCmd.PersistentFlags().BoolVar(
		&noPrettyOutput, "no-pretty", false, "Disable pretty print JSON outputs",
	)
	RootCmd.AddCommand(&cobra.Command{
		Use:   "cookbook",
		Short: "Example usages of the most important, frequently used commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(cookbookDoc)
			return nil
		},
	})
	RootCmd.PersistentFlags().BoolVarP(&useLedgerWallet, "ledger", "e", false, "Use ledger hardware wallet")
	RootCmd.PersistentFlags().StringVar(&givenFilePath, "file", "", "Path to file for given command when applicable")
	RootCmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: fmt.Sprintf("Generate docs to a local %s directory", hmyDocsDir),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			docDir := path.Join(cwd, hmyDocsDir)
			os.Mkdir(docDir, 0700)
			doc.GenMarkdownTree(RootCmd, docDir)
			return nil
		},
	})
}

var (
	// VersionWrapDump meant to be set from main.go
	VersionWrapDump = ""
	cookbook        = color.GreenString("hmy cookbook")
)

// Execute kicks off the hmy CLI
func Execute() {
	RootCmd.SilenceErrors = true
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(errors.Wrapf(err, "commit: %s, error", VersionWrapDump).Error())
		fmt.Println("check " + cookbook + " for valid examples or try adding a `--help` flag")
		os.Exit(1)
	}
}

func validateAddress(cmd *cobra.Command, args []string) error {
	// Check if input valid one address
	address := oneAddress{}
	if err := address.Set(args[0]); err != nil {
		// Check if input is valid account name
		if acc, err := store.AddressFromAccountName(args[0]); err == nil {
			addr = oneAddress{acc}
			return nil
		}
		return fmt.Errorf("Invalid one address/Invalid account name: %s", args[0])
	}
	addr = address
	return nil
}
