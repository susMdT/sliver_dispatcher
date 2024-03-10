package utils

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sliver-dispatch/globals"
	"sort"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

const (
	EXT_EXEC_SUCCESS = 1
	EXT_EXEC_ERR     = 2
	EXT_ERR_OTHER    = 3
)

func WinExtensionExec(args []string, ExtCfg globals.ExtensionCfg, session *clientpb.Session, rpc rpcpb.SliverRPCClient) int {

	var resp_l *sliverpb.ListExtensions
	var resp_r *sliverpb.RegisterExtension
	var resp_c *sliverpb.CallExtension
	var extArgs []byte
	var err error

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
		Eprint("Error checking %s is loaded into the session: "+err.Error(), ExtCfg.Command_Name)
		return EXT_ERR_OTHER
	}

	if !slices.Contains(resp_l.Names, ExtCfg.Command_Name) {
		if session.Arch == "amd64" {

			b_ext, err := os.ReadFile(
				"extensions/" +
					ExtCfg.Command_Name + "/" +
					ExtCfg.Files[sort.Search(len(ExtCfg.Files), func(i int) bool { return ExtCfg.Files[i].Arch == session.Arch })].Path)

			if err != nil {
				Eprint("Error reading file: %s", err.Error())
				return EXT_ERR_OTHER
			}
			Dprint("Adding %s extension to this session", ExtCfg.Command_Name)
			resp_r, err = rpc.RegisterExtension(
				context.Background(),
				&sliverpb.RegisterExtensionReq{
					Name: ExtCfg.Command_Name,
					OS:   "windows",
					Data: b_ext,
					Request: &commonpb.Request{
						Async:     false,
						SessionID: session.ID,
					},
				},
			)
			if resp_r != nil {
				if resp_r.Response != nil {
					if resp_r.Response.Err != "" {
						Eprint(resp_r.Response.Err)
						return EXT_ERR_OTHER
					}
				}
			}
			if err != nil {
				Eprint("Error loading %s extension to the session: "+err.Error(), ExtCfg.Command_Name)
				return EXT_ERR_OTHER
			}
		} else {
			Eprint("Extension execution not supported for x86!")
			return EXT_ERR_OTHER
		}
	}

	extArgs, err = GetExtArgs(
		args,
		"extensions/"+ExtCfg.Command_Name+"/"+
			ExtCfg.Files[sort.Search(len(ExtCfg.Files), func(i int) bool { return ExtCfg.Files[i].Arch == session.Arch })].Path,
		ExtCfg,
	)
	fmt.Println(args)
	fmt.Println(extArgs)
	if err != nil {
		Eprint("Error parsing extension arguments: " + err.Error())
		return EXT_ERR_OTHER
	}

	resp_c, err = rpc.CallExtension(
		context.Background(),
		&sliverpb.CallExtensionReq{
			Name:   ExtCfg.Command_Name,
			Export: ExtCfg.Entrypoint,
			Args:   extArgs,
			Request: &commonpb.Request{
				Async:     false,
				SessionID: session.ID,
			},
		},
	)

	if resp_c != nil {
		if resp_c.Output != nil {
			Iprint("Response: " + string(resp_c.Output))
		}
		if resp_c.Response != nil {
			if resp_c.Response.Err != "" {
				Eprint("Error: " + resp_c.Response.Err)
			}
		}
	}
	if err != nil {
		Eprint("Error: " + err.Error())
	}

	if err != nil {
		return EXT_EXEC_ERR
	}
	if resp_c.Response != nil {
		if resp_c.Response.Err != "" {
			return EXT_EXEC_ERR
		}
	}
	return EXT_EXEC_SUCCESS

}
