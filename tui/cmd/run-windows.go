package cmd

import (
	"github.com/spf13/cobra"
)

var Run_all_windows = &cobra.Command{
	Use:                "run_all_windows [ module ]",
	Short:              "Run a module across all windows sessions",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
