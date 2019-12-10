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
		RunE: func(cmd *cobra.Command, args []string) error {
			RootCmd.GenBashCompletion(os.Stdout)
			return nil
		},
	}

	cmdStaking.AddCommand(stakingSubCommands()...)
	RootCmd.AddCommand(cmdStaking)
}
