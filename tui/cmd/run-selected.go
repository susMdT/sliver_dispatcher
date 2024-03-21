package cmd

import "github.com/spf13/cobra"

var Run_all_selected = &cobra.Command{
	Use:                "run_all_selected [ module ]",
	Short:              "Run a module across all selected sessions",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
