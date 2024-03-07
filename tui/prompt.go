package tui

import (
	"fmt"
	"os"
	"sliver-dispatch/client_cmds"
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"strconv"
	"strings"

	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "get_sessions", Description: "List active sliver sessions"},
		{Text: "get_configs", Description: "List prior sliver build configurations and server configs"},
		{Text: "select_sessions", Description: "Select sliver sessions to run commands on"},
		{Text: "run_all_selected", Description: "Run a command across all selected sessions"},
		{Text: "run_all_windows", Description: "Run a command across all windows sessions"},
		{Text: "run_all_linux", Description: "Run a command across all linux sessions"},
		{Text: "help", Description: "Help"},
		{Text: "toggle_debug", Description: "Toggle debug mode"},
		{Text: "exit", Description: "Exit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Main(rpc rpcpb.SliverRPCClient) {

	for {
		t := prompt.Input("[>] ", completer)
		input := utils.SplitArguments(t)
		if len(input) > 0 {
			utils.Dprint("Input: %s", t)
			switch strings.ToLower(input[0]) {
			case "exit":
				os.Exit(0)
			case "get_sessions":
				utils.UpdateSessions(rpc)
				client_cmds.GetSessions()
			case "get_configs":
				client_cmds.GetConfigs(rpc)
			case "run_all_windows":
				client_cmds.RunMass(rpc, input[1:]...)
			case "debug":
				globals.DebugMode = !globals.DebugMode
				fmt.Println("Debug mode set to " + strconv.FormatBool(globals.DebugMode))
			default:
				fmt.Println("Unknown command: " + t)
			}
		}
		utils.UpdateSessions(rpc)

	}
}
