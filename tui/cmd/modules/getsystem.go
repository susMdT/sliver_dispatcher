package modules

import (
	"sliver-dispatch/client_cmds/dispatch"
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"strings"

	"github.com/spf13/cobra"
)

var GetSystem = &cobra.Command{
	Use:   "getsystem [ /path/to/shc ] [ process_to_spawn ]",
	Short: "Spawn a system process and inject shellcode",
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
			dispatch.GetSystem(globals.Rpc, args[0:]...)
		} else {
			utils.Eprint("Cannot run getsystem on linux!")
		}
	},
}
