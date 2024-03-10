package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sliver-dispatch/globals"
	"sort"
	"strconv"
	"strings"

	"github.com/bishopfox/sliver/client/core"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
)

const (
	BOF_EXEC_SUCCESS = 1
	BOF_EXEC_ERR     = 2
	BOF_ERR_OTHER    = 3
)

func GetBOFArgs(args []string, binPath string, Ext globals.ExtensionCfg) ([]byte, error) {
	var extensionArgs []byte
	binData, err := os.ReadFile(binPath)
	if err != nil {
		return nil, err
	}

	// Now build the extension's argument buffer
	extensionArgsBuffer := core.BOFArgsBuffer{
		Buffer: new(bytes.Buffer),
	}
	err = extensionArgsBuffer.AddString(Ext.Entrypoint)
	if err != nil {
		return nil, err
	}
	Dprint("Added entrypoint argument: %s", Ext.Entrypoint)

	err = extensionArgsBuffer.AddData(binData)
	if err != nil {
		return nil, err
	}
	Dprint("Added binary data, length: %d", len(binData))

	parsedArgs, err := GetExtArgs(args, binPath, Ext)
	if err != nil {
		return nil, err
	}
	err = extensionArgsBuffer.AddData(parsedArgs)
	if err != nil {
		return nil, err
	}
	Dprint("Added parsed arguments, length: %d", len(parsedArgs))

	extensionArgs, err = extensionArgsBuffer.GetBuffer()
	if err != nil {
		return nil, err
	}
	return extensionArgs, nil
}

func GetExtArgs(args []string, binPath string, Ext globals.ExtensionCfg) ([]byte, error) {
	var err error
	argsBuffer := core.BOFArgsBuffer{
		Buffer: new(bytes.Buffer),
	}

	// Parse BOF arguments from grumble
	missingRequiredArgs := make([]string, 0)

	// If we have an extension that expects a single string, but more than one has been parsed, combine them
	if len(Ext.Arguments) == 1 && strings.Contains(Ext.Arguments[0].Type, "string") {
		// The loop below will only read the first element of args because ext.Arguments is 1
		args[0] = strings.Join(args, " ")
	}

	for _, arg := range Ext.Arguments {
		// If we don't have any positional words left to consume,
		// add the remaining required extension arguments in the
		// error message.
		if len(args) == 0 {
			if !arg.Optional {
				missingRequiredArgs = append(missingRequiredArgs, "`"+arg.Name+"`")
			}
			continue
		}

		// Else pop a word from the list
		word := args[0]
		args = args[1:]
		switch arg.Type {
		case "integer":
			fallthrough
		case "int":
			val, err := strconv.Atoi(word)
			if err != nil {
				return nil, err
			}
			err = argsBuffer.AddInt(uint32(val))
			if err != nil {
				return nil, err
			}
		case "short":
			val, err := strconv.Atoi(word)
			if err != nil {
				return nil, err
			}
			err = argsBuffer.AddShort(uint16(val))
			if err != nil {
				return nil, err
			}
		case "string":
			err = argsBuffer.AddString(word)
			if err != nil {
				return nil, err
			}
		case "wstring":
			err = argsBuffer.AddWString(word)
			if err != nil {
				return nil, err
			}
		// Adding support for filepaths so we can
		// send binary data like shellcodes to BOFs
		case "file":
			data, err := os.ReadFile(word)
			if err != nil {
				return nil, err
			}
			err = argsBuffer.AddData(data)
			if err != nil {
				return nil, err
			}
		}
	}

	// Return if we have missing required arguments
	if len(missingRequiredArgs) > 0 {
		return nil, fmt.Errorf("required arguments %s were not provided", strings.Join(missingRequiredArgs, ", "))
	}

	parsedArgs, err := argsBuffer.GetBuffer()
	if err != nil {
		return nil, err
	}

	return parsedArgs, nil
}

func ParseExtCfg(path string) (globals.ExtensionCfg, error) {

	var CfgParsed globals.ExtensionCfg

	CfgData, err := os.ReadFile(path)
	if err != nil {
		return globals.ExtensionCfg{}, err
	}

	err = json.Unmarshal(CfgData, &CfgParsed)
	if err != nil {
		return globals.ExtensionCfg{}, err
	}

	return CfgParsed, nil
}

func BofExec(bofname string, args []string, session *clientpb.Session, rpc rpcpb.SliverRPCClient) int {

	b_coff64, err := os.ReadFile("extensions/coff-loader/COFFLoader.x64.dll")
	if err != nil {
		Eprint("Error reading file: %s", err.Error())
		return BOF_ERR_OTHER
	}

	b_coff32, err := os.ReadFile("extensions/coff-loader/COFFLoader.x86.dll")
	if err != nil {
		Eprint("Error reading file: %s", err.Error())
		return BOF_ERR_OTHER
	}

	var resp_l *sliverpb.ListExtensions
	var resp_r *sliverpb.RegisterExtension
	var resp_c *sliverpb.CallExtension
	var ldrCfg globals.ExtensionCfg
	var bofCfg globals.ExtensionCfg
	var extArgs []byte
	if !session.IsDead && session.OS == "windows" {

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
			Eprint("Error checking coff-loader is loaded into the session: " + err.Error())
			return BOF_ERR_OTHER
		}

		if !slices.Contains(resp_l.Names, "coff-loader") {
			Dprint("Registering coff-loader extension")
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
				if resp_r.Response != nil {
					if resp_r.Response.Err != "" {
						Eprint(resp_r.Response.Err)
						return BOF_ERR_OTHER
					}
				}
			}
			if err != nil {
				Eprint("Error loading coff-loader extension to the session: " + err.Error())
				return BOF_ERR_OTHER
			}
		}
		ldrCfg, err = ParseExtCfg("extensions/coff-loader/extension.json")
		if err != nil {
			Eprint("Error parsing the coff-loader extension configuration: " + err.Error())
			return BOF_ERR_OTHER
		}

		bofCfg, err = ParseExtCfg("extensions/" + bofname + "/extension.json")
		if err != nil {
			Eprint("Error parsing %s extension configuration: "+err.Error(), bofname)
			return BOF_ERR_OTHER
		}
		extArgs, err = GetBOFArgs(
			args,
			"extensions/"+bofCfg.Command_Name+"/"+bofCfg.Files[sort.Search(len(bofCfg.Files), func(i int) bool { return bofCfg.Files[i].Arch == session.Arch })].Path,
			bofCfg,
		)
		if err != nil {
			Eprint("Error parsing extension arguments: " + err.Error())
			return BOF_ERR_OTHER
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
			Iprint("Response: " + string(resp_c.Output))
		}
		if err != nil {
			Eprint("Error: " + err.Error())
		}
		if resp_c.Response != nil {
			if resp_c.Response.Err != "" {
				Eprint("Error: " + resp_c.Response.Err)
			}
		}

		if err != nil {
			return BOF_EXEC_ERR
		}
		if resp_c.Response != nil {
			if resp_c.Response.Err != "" {
				return BOF_EXEC_ERR
			}
		}
		return BOF_EXEC_SUCCESS
	}
	return BOF_ERR_OTHER
}
