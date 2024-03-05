package client_cmds

import (
	"sliver-dispatch/client_cmds/dispatch"
	"strings"

	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

func RunMass(rpc rpcpb.SliverRPCClient, args ...string) {

	switch strings.ToLower(args[0]) {
	case "execute":
		dispatch.Execute(rpc, args[1:]...)
	}
}
