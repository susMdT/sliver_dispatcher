package dispatch

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

func Upload(rpc rpcpb.SliverRPCClient, args ...string) {

	if len(args) < 2 {
		fmt.Println("Need the full path to a file to upload and the destination path!")
		return
	}
	var sessions *clientpb.Sessions

	sessions, err := rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var session *clientpb.Session
	for _, session = range sessions.GetSessions() {
		if !session.IsDead {
			_, err = rpc.Upload(
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
			if err != nil {
				fmt.Println(fmt.Sprintf("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s",
					strings.Split(session.ID, "-")[0],
					session.Hostname,
					strings.Split(session.RemoteAddress, ":")[0],
					session.Username))
				fmt.Println("Error: " + err.Error())
			}
		}
	}
}
