package cmd

import (
	"sliver-dispatch/client_cmds"

	"github.com/spf13/cobra"
)

var GetSessions = &cobra.Command{
	Use:   "get_sessions",
	Short: "Get all active sessions",
	Run: func(cmd *cobra.Command, args []string) {
		client_cmds.GetSessions()
	},
}
