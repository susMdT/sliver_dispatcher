package cmd

import (
        "sliver-dispatch/tui/cmd/modules"

        "github.com/c-bata/go-prompt"
        "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
        Use:           "",
        SilenceUsage:  true, // Only print usage when defined in command.
        SilenceErrors: true,
}

var mainCmds = []prompt.Suggest{
        {Text: "get_sessions", Description: "List active sliver sessions"},
        {Text: "get_configs", Description: "List prior sliver build configurations and server configs"},
        {Text: "select_sessions", Description: "Select sliver sessions to run commands on"},
        {Text: "run_all_windows", Description: "Run a command across all windows sessions"},
        {Text: "run_all_linux", Description: "Run a command across all linux sessions"},
        {Text: "help", Description: "Help"},
        {Text: "toggle_debug", Description: "Toggle debug mode"},
        {Text: "exit", Description: "Exit"},
}
var GetMainCmd = func(main_cmd string) []prompt.Suggest {

        cmdExists := false
        for _, mCmd := range mainCmds {
                if mCmd.Text == main_cmd {
                        cmdExists = true
                        break
                }
        }

        if cmdExists {
                return mainCmds
        }
        return nil
}

func init() {

        RootCmd.AddCommand(Debug)
        RootCmd.AddCommand(GetConfigs)
        RootCmd.AddCommand(GetSessions)
        RootCmd.AddCommand(SelectSessions)

        RootCmd.AddCommand(Run_all_windows)
        Run_all_windows.AddCommand(modules.GetExecuteInst())
        Run_all_windows.AddCommand(modules.GetUploadInst())
        Run_all_windows.AddCommand(modules.GetSystem)
        Run_all_windows.AddCommand(modules.Nosferatu)
        Run_all_windows.AddCommand(modules.Shinject)
        Run_all_windows.AddCommand(modules.GetScriptInst())

        RootCmd.AddCommand(Run_all_linux)
        Run_all_linux.AddCommand(modules.GetExecuteInst())
        Run_all_linux.AddCommand(modules.GetUploadInst())
        Run_all_linux.AddCommand(modules.GetScriptInst())

        RootCmd.AddCommand(Run_all_selected)
        Run_all_selected.AddCommand(modules.GetExecuteInst())
        Run_all_selected.AddCommand(modules.GetUploadInst())
}

