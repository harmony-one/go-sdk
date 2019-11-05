package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"

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
			fmt.Print(string(asJSON))
			return nil
		}
		fmt.Print(common.JSONPrettyFormat(string(asJSON)))
		return nil
	}
	// RootCmd is single entry point of the CLI
	RootCmd = &cobra.Command{
		Use:          "hmy",
		Short:        "Harmony blockchain",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				common.EnableAllVerbose()
			}
		},
		Long: fmt.Sprintf(`
CLI interface to the Harmony blockchain

%s`, g("Invoke 'hmy cookbook' for examples of the most common, important usages")),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	vS := "dump out debug information, same as env var HMY_ALL_DEBUG=true"
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, vS)
	RootCmd.PersistentFlags().StringVarP(&node, "node", "n", defaultNodeAddr, "<host>")
	RootCmd.PersistentFlags().BoolVar(&noLatest, "no-latest", false, "Do not add 'latest' to RPC params")
	RootCmd.PersistentFlags().BoolVar(&noPrettyOutput, "no-pretty", false, "Disable pretty print JSON outputs")
	RootCmd.AddCommand(&cobra.Command{
		Use:   "cookbook",
		Short: "Example usages of the most important, frequently used commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(cookbookDoc)
		},
	})
	RootCmd.PersistentFlags().BoolVarP(&useLedgerWallet, "ledger", "e", false, "Use ledger hardware wallet")
	RootCmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: fmt.Sprintf("Generate docs to a local %s directory", hmyDocsDir),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, _ := os.Getwd()
			docDir := path.Join(cwd, hmyDocsDir)
			os.Mkdir(docDir, 0700)
			doc.GenMarkdownTree(RootCmd, docDir)
		},
	})
}

// Execute kicks off the hmy CLI
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
