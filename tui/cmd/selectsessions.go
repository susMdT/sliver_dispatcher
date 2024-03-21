package cmd

import (
        "sliver-dispatch/client_cmds"

        "github.com/spf13/cobra"
)

var SelectSessions = &cobra.Command{
        Use:   "select_sessions",
        Short: "Select active sessions",
        Run: func(cmd *cobra.Command, args []string) {
                client_cmds.SelectSessions()
        },
}
