package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

var (
	useLatestInParamsForRPC   bool
	prettyPrintJSONOutput     bool
	useOneAddressInsteadOfHex bool
	node                      string
	keyStoreDir               string
	request                   = func(method rpc.RPCMethod, params []interface{}) {
		if useLatestInParamsForRPC {
			params = append(params, "latest")
		}
		asJSON, _ := json.Marshal(rpc.RPCRequest(method, node, params))
		if prettyPrintJSONOutput {
			fmt.Print(jsonPrettyPrint(string(asJSON)))
			return
		}
		fmt.Print(string(asJSON))
	}
	RootCmd = &cobra.Command{
		Use:   "hmy_cli",
		Short: "Harmony blockchain",
		Long: `
CLI interface to the Harmony blockchain
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

const (
	HMY_CLI_DOCS_DIR  = "hmy_cli-docs"
	DEFAULT_NODE_ADDR = "http://localhost:9500"
)

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&node,
		"node",
		"n",
		DEFAULT_NODE_ADDR,
		"<host>:<port>",
	)
	RootCmd.PersistentFlags().BoolVarP(&useLatestInParamsForRPC, "latest", "l", false, "Add 'latest' to RPC params")
	RootCmd.PersistentFlags().BoolVarP(&prettyPrintJSONOutput, "pretty", "p", false, "pretty print JSON outputs")
	RootCmd.PersistentFlags().BoolVarP(&useOneAddressInsteadOfHex, "one-address", "o", false, "use one address for RPC calls")
	RootCmd.PersistentFlags().StringVar(&keyStoreDir, "key-store-dir", "k", "What directory to use as the keystore")
	RootCmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: fmt.Sprintf("Generate docs to a local %s directory", HMY_CLI_DOCS_DIR),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, _ := os.Getwd()
			docDir := path.Join(cwd, HMY_CLI_DOCS_DIR)
			os.Mkdir(docDir, 0700)
			doc.GenMarkdownTree(RootCmd, docDir)
		},
	})
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
