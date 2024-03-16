package cmd

import "github.com/spf13/cobra"

var Run_all_linux = &cobra.Command{
	Use:                "run_all_linux [ module ]",
	Short:              "Run a module across all linux sessions",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
