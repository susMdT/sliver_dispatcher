package utils

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sliver-dispatch/globals"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	. "github.com/logrusorgru/aurora"
)

func UpdateSessions(rpc rpcpb.SliverRPCClient) {
	var sessions *clientpb.Sessions
	sessions, err := rpc.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	var session *clientpb.Session
	globals.ActiveSessions = globals.ActiveSessions[:0]
	for _, session = range sessions.Sessions {
		if !session.IsDead {
			globals.ActiveSessions = append(globals.ActiveSessions, globals.Interactive_Session{
				Host:     session.Hostname,
				OS:       session.OS,
				ID:       strings.Split(session.ID, "-")[0],
				PID:      int(session.PID),
				Address:  strings.Split(session.RemoteAddress, ":")[0],
				Process:  session.Filename,
				Username: session.Username,
			})
		}
	}
}

func SplitArguments(userInput string) []string {
	// Define a regular expression pattern to match quoted substrings or non-quoted substrings
	pattern := `"(?:\\.|[^"\\])*"|\S+`

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)

	// Use FindAllString() to find all matches in the user input
	matches := re.FindAllString(userInput, -1)

	// Remove the quotes unless they are escaped
	var arguments []string
	for _, match := range matches {
		arguments = append(arguments, strings.Trim(match, `"`))
	}

	return arguments
}

func Dprint(str string, args ...interface{}) {
	if globals.DebugMode {
		fmt.Printf("%s %s\n", Cyan("[*]"), fmt.Sprintf(str, args...))
	}
}

func Eprint(str string, args ...interface{}) {
	fmt.Printf("%s %s\n", BrightRed("[!]"), fmt.Sprintf(str, args...))
}

func Iprint(str string, args ...interface{}) {
	fmt.Printf("%s %s\n", BrightGreen("[+]"), fmt.Sprintf(str, args...))
}
