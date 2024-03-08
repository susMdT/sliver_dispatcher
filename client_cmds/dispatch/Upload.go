package dispatch

import (
	"context"
	"fmt"
	"os"
	"sliver-dispatch/utils"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

func Upload(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 2 {
		utils.Eprint("Need the full path to a file to upload and the destination path!")
		return
	}
	var sessions *clientpb.Sessions

	sessions, err := rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error listing sessions: %s", err.Error())
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		utils.Eprint("Error reading file: %s", err.Error())
		return
	}

	var session *clientpb.Session
	var upload_rsp *sliverpb.Upload
	for _, session = range sessions.GetSessions() {
		if !session.IsDead {
			upload_rsp, err = rpc.Upload(
				context.Background(), &sliverpb.UploadReq{
					Path:    args[1],
					Data:    data,
					IsIOC:   false,
					Encoder: "",
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
				utils.Eprint("Error: " + err.Error())
				if upload_rsp != nil {
					if upload_rsp.Response != nil && upload_rsp.Response.Err != "" {
						utils.Eprint("Another error: " + upload_rsp.Response.Err)
					}
				}
			}
		}
	}
}
