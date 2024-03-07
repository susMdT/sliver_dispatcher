package client_cmds

import (
	"fmt"
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
		fmt.Println("Missing arguments!")
		return
	}
	switch strings.ToLower(args[0]) {
	case "execute":
		dispatch.Execute(rpc, args[1:]...)
	case "upload":
		dispatch.Upload(rpc, args[1:]...)
	}
}
