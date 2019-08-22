package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hmy_cli",
	Short: "Harmony blockchain",
	Long: `
CLI interface to the Harmony blockchain
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	prettyPrintJSONOutput bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prettyPrintJSONOutput, "pretty", "p", false, "pretty print JSON outputs")

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {

			os.Exit(0)
		},
	}

	rootCmd.AddCommand(cmdVersion)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
