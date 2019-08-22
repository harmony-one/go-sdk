package main

import (
	"fmt"
	"os"

	"github.com/harmony-one/go-sdk/client/cmd"
	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	builtAt string
	builtBy string
)

func main() {
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s %s %s %s", version, commit, builtAt, builtBy)
			os.Exit(0)
		},
	})
	cmd.Execute()
}
