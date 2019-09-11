package main

import (
	"fmt"
	"os"
	"path"

	"github.com/harmony-one/go-sdk/client/cmd"
	_ "github.com/harmony-one/go-sdk/pkg/store"
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
			fmt.Fprintf(os.Stderr,
				"Harmony (C) 2019. %v, version %v-%v (%v %v)\n",
				path.Base(os.Args[0]), version, commit, builtBy, builtAt)
			os.Exit(0)
		},
	})
	cmd.Execute()
}
