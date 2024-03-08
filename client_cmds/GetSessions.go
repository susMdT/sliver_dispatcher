package client_cmds

import (
	"sliver-dispatch/globals"
	"sliver-dispatch/utils"
)

func GetSessions() {
	utils.Iprint("[Alive Sessions]")

	var session globals.Interactive_Session
	for _, session = range globals.ActiveSessions {
		utils.Iprint("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s | Process: %-10s | PID: %d",
			session.ID,
			session.Host,
			session.Address,
			session.Username,
			session.Process,
			session.PID)
	}

	utils.Iprint("[Selected Sessions]")
	for _, session = range globals.Selected_Sessions {
		utils.Iprint("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s | Process: %-10s | PID: %d",
			session.ID,
			session.Host,
			session.Address,
			session.Username,
			session.Process,
			session.PID)
	}

}
