package client_cmds

import (
	"fmt"
	"sliver-dispatch/globals"
)

func GetSessions() {
	fmt.Println("[Alive Sessions]")

	var session globals.Interactive_Session
	for _, session = range globals.ActiveSessions {
		fmt.Printf(fmt.Sprintf("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s | Process: %-10s | PID: %d\n",
			session.ID,
			session.Host,
			session.Address,
			session.Username,
			session.Process,
			session.PID))
	}

	fmt.Println("[Selected Sessions]")
	for _, session = range globals.Selected_Sessions {
		fmt.Printf(fmt.Sprintf("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s | Process: %-10s | PID: %d\n",
			session.ID,
			session.Host,
			session.Address,
			session.Username,
			session.Process,
			session.PID))
	}

}
