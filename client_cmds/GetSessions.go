package client_cmds

import (
    "sliver-dispatch/globals"
    "sliver-dispatch/utils"
)

func GetSessions() {
    utils.Iprint("[Alive Sessions]")
    globals.L_rpc.Lock()
    for _, session := range globals.ActiveSessions {
        utils.Iprint("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s | Process: %-10s | PID: %d",
            session.ID,
            session.Host,
            session.Address,
            session.Username,
            session.Process,
            session.PID)
    }
    globals.L_rpc.Unlock()

    utils.Iprint("[Selected Sessions]")
    // Locking here again, ensure you handle locking according to your application's concurrency model
    globals.L_rpc.Lock()
    defer globals.L_rpc.Unlock() // Using defer to ensure the lock is released
    for _, selectedID := range globals.Selected_Sessions {
        for _, session := range globals.ActiveSessions {
            if session.ID == selectedID {
                utils.Iprint("ID: %-10s | Host: %-20s | Address: %-15s | Username: %-10s | Process: %-10s | PID: %d",
                    session.ID,
                    session.Host,
                    session.Address,
                    session.Username,
                    session.Process,
                    session.PID)
                break // Found the matching session, no need to check further
            }
        }
    }
}

