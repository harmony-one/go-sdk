package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cmdStaking := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion, run:

    . <(hmy completion)

Add the line to your ~/.bashrc to enable completiony for each bash session.
`,
		Run: func(cmd *cobra.Command, args []string) {
			RootCmd.GenBashCompletion(os.Stdout)
		},
	}

	cmdStaking.AddCommand(stakingSubCommands()...)
	RootCmd.AddCommand(cmdStaking)
}
