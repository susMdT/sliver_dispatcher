package cmd

import (
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"strconv"

	"github.com/spf13/cobra"
)

var Debug = &cobra.Command{
	Use:   "debug",
	Short: "Enable debug mode",
	Run: func(cmd *cobra.Command, args []string) {

		globals.DebugMode = !globals.DebugMode
		utils.Iprint("Set debug mode to %s", strconv.FormatBool(globals.DebugMode))
	},
}
