package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sliver-dispatch/globals"
	"strconv"
	"strings"

	"github.com/bishopfox/sliver/client/core"
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
