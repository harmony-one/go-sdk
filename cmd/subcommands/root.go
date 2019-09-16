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
	useLedgerWallet bool
	noLatest        bool
	noPrettyOutput  bool
	node            string
	keyStoreDir     string
	request         = func(method rpc.RPCMethod, params []interface{}) {
		if !noLatest {
			params = append(params, "latest")
		}
		success, failure := rpc.Request(method, node, params)
		if failure != nil {
			fmt.Println(failure)
			os.Exit(-1)
		}
		asJSON, _ := json.Marshal(success)
		if noPrettyOutput {
			fmt.Print(string(asJSON))
			return
		}
		fmt.Print(common.JSONPrettyFormat(string(asJSON)))
	}
	RootCmd = &cobra.Command{
		Use:          "hmy",
		Short:        "Harmony blockchain",
		SilenceUsage: true,
		Long: fmt.Sprintf(`
CLI interface to the Harmony blockchain

%s`, g("Invoke 'hmy cookbook' for examples of the most common, important usages")),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
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

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
