package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	prettyPrintJSONOutput bool
	node                  string
	RootCmd               = &cobra.Command{
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
	HMY_CLI_DOCS_DIR = "hmy_cli-docs"
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&prettyPrintJSONOutput, "pretty", "p", false, "pretty print JSON outputs")
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
