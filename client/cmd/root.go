package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

func init() {
	RootCmd.PersistentFlags().BoolVarP(&prettyPrintJSONOutput, "pretty", "p", false, "pretty print JSON outputs")

}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
