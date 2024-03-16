package cmd

import (
	"sliver-dispatch/client_cmds"
	"sliver-dispatch/globals"

	"github.com/spf13/cobra"
)

var GetConfigs = &cobra.Command{
	Use:   "get_configs",
	Short: "Get implant configs and profiles",
	Run: func(cmd *cobra.Command, args []string) {
		client_cmds.GetConfigs(globals.Rpc)
	},
}
