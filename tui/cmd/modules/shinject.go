package modules

import (
	"sliver-dispatch/client_cmds/dispatch"
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"strings"

	"github.com/spf13/cobra"
)

var Shinject = &cobra.Command{
	Use:   "shinject [ source_path ] [ process_name ]",
	Short: "Inject shellcode into a given process name",
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
			dispatch.Shinject(globals.Rpc, args[0:]...)
		} else if 
		   cmd.Parent().Name() == "run_all_selected" {
			dispatch.ShinjectOnSelectedSessions(globals.Rpc, args[0:]...)
		} else {
			utils.Eprint("Cannot run shinject on linux!")
		}

	},
}
