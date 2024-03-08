package dispatch

import (
	"context"
	"fmt"
	"log"
	"sliver-dispatch/utils"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

// https://github.com/BishopFox/sliver/blob/c8a7948671eafba4d6f871c2f2b46b900202699d/client/command/exec/execute.go#L71
// https://github.com/BishopFox/sliver/blob/c8a7948671eafba4d6f871c2f2b46b900202699d/client/console/console.go#L770
func Execute(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 1 {
		fmt.Println("Need the name of an executable to run!")
		return
	}
	var sessions *clientpb.Sessions

	sessions, err := rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	var session *clientpb.Session
	var exec *sliverpb.Execute
	for _, session = range sessions.GetSessions() {
		if !session.IsDead {
			exec, err = rpc.Execute(
				context.Background(),
				&sliverpb.ExecuteReq{
					Path:   args[0],
					Args:   args[1:],
					Output: true,
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				})

			utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
				strings.Split(session.ID, "-")[0],
				session.Hostname,
				strings.Split(session.RemoteAddress, ":")[0],
				session.Username))
			if err != nil {
				utils.Eprint("Error: %s", err.Error())
			}
			utils.Iprint(string(exec.Stdout) + string(exec.Stderr))
		}
	}
}
