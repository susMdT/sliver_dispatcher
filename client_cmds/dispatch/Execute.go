package dispatch

import (
	"context"
	"fmt"
	"sliver-dispatch/utils"
	"strings"
	"sliver-dispatch/globals"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

// https://github.com/BishopFox/sliver/blob/c8a7948671eafba4d6f871c2f2b46b900202699d/client/command/exec/execute.go#L71
// https://github.com/BishopFox/sliver/blob/c8a7948671eafba4d6f871c2f2b46b900202699d/client/console/console.go#L770
func Execute(rpc rpcpb.SliverRPCClient, target_os string, args ...string) {

	if len(args) < 1 {
		utils.Eprint("Need the name of an executable to run!")
		return
	}
	var sessions *clientpb.Sessions

	sessions, err := rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error reading file: %s", err.Error())
		return
	}

	var session *clientpb.Session
	var exec *sliverpb.Execute
	for _, session = range sessions.GetSessions() {
		if !session.IsDead && session.OS == target_os {
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
			if exec != nil {
				utils.Iprint(string(exec.Stdout) + string(exec.Stderr))
			}
		}
	}
}

func ExecuteOnSelectedSessions(rpc rpcpb.SliverRPCClient, args ...string) {
    if len(args) < 1 {
        utils.Eprint("Need the name of an executable to run!")
        return
    }

    sessions, err := rpc.GetSessions(context.Background(), &commonpb.Empty{})
    if err != nil {
        utils.Eprint("Error retrieving sessions: %s", err.Error())
        return
    }

    // Create a map for efficient lookup of selected session IDs (first part only)
    selectedSessionsMap := make(map[string]struct{})
    for _, id := range globals.Selected_Sessions {
        idParts := strings.SplitN(id, "-", 2)
        if len(idParts) > 0 {
            selectedSessionsMap[idParts[0]] = struct{}{}
        }
    }

    for _, session := range sessions.GetSessions() {
        idParts := strings.SplitN(session.ID, "-", 2)
        firstPartID := ""
        if len(idParts) > 0 {
            firstPartID = idParts[0]
        }
        
        if _, isSelected := selectedSessionsMap[firstPartID]; isSelected && !session.IsDead {
            exec, err := rpc.Execute(
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

            // Log session info and command output
            utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
                strings.Split(session.ID, "-")[0],
                session.Hostname,
                strings.Split(session.RemoteAddress, ":")[0],
                session.Username))
            if err != nil {
                utils.Eprint("Error executing command on session %s: %s", session.ID, err.Error())
                continue
            }
            if exec != nil {
                utils.Iprint(string(exec.Stdout) + string(exec.Stderr))
            }
        }
    }
}
