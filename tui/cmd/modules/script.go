package modules

import (
	"sliver-dispatch/client_cmds/dispatch"
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"strings"

	"github.com/spf13/cobra"
)

func GetScriptInst() *cobra.Command {
	return &cobra.Command{
		Use:   "script [ source_path ]",
		Short: "Run a script across sessions",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			args = utils.SplitArguments(strings.Join(args, " "))
			for i, arg := range args {
				utils.Dprint("Arg %d: %s", i, arg)
			}
			help, _ := cmd.Flags().GetBool("help")

			if len(args) == 0 || help {
				cmd.Usage()
				return
			}
			if cmd.Parent().Name() == "run_all_windows" {

				dispatch.Script(globals.Rpc, "windows", args[0:]...)
			} else if cmd.Parent().Name() == "run_all_linux" {

				dispatch.Script(globals.Rpc, "linux", args[0:]...)
			}
		},
	}
}
