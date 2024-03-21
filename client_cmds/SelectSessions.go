package client_cmds

import (
    "sort"
    "sliver-dispatch/globals"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

type sessionSelector struct {
    app                *tview.Application
    list               *tview.List
    selectedSessionIDs map[string]struct{}
    lastIndex          int
}

func newSessionSelector() *sessionSelector {
    ss := &sessionSelector{
        app:                tview.NewApplication(),
        list:               tview.NewList().ShowSecondaryText(false),
        selectedSessionIDs: make(map[string]struct{}),
    }
    ss.list.SetHighlightFullLine(true)
    return ss
}

func (ss *sessionSelector) updateList() {
    ss.list.Clear()

    globals.L_rpc.Lock()
    sessionsCopy := make([]globals.Interactive_Session, len(globals.ActiveSessions))
    copy(sessionsCopy, globals.ActiveSessions)
    globals.L_rpc.Unlock()

    sort.SliceStable(sessionsCopy, func(i, j int) bool {
        return sessionsCopy[i].Address < sessionsCopy[j].Address
    })

    // Update the list items and their selection callbacks
    for _, session := range sessionsCopy {
        sessionID := session.ID // Capture the session ID for use in the callback
        _, isSelected := ss.selectedSessionIDs[sessionID]
        prefix := " "
        if isSelected {
            prefix = "[*] "
        }
        itemText := prefix + "ID: " + sessionID + " | Host: " + session.Host +
                    " | Address: " + session.Address + " | Username: " + session.Username +
                    " | Process: " + session.Process

        ss.list.AddItem(itemText, "", 0, func() {
            if _, ok := ss.selectedSessionIDs[sessionID]; ok {
                delete(ss.selectedSessionIDs, sessionID) // Deselect
            } else {
                ss.selectedSessionIDs[sessionID] = struct{}{} // Select
            }
            ss.lastIndex = ss.list.GetCurrentItem() // Store the current index
            ss.updateList()                          // Refresh the list to reflect the new selection state
            ss.list.SetCurrentItem(ss.lastIndex)    // Set the cursor back to the last selected item
        })
    }
    ss.list.SetCurrentItem(ss.lastIndex) // Set the cursor to the last selected item
}

func (ss *sessionSelector) run() {
    ss.updateList()
    ss.list.SetBorder(true).SetTitle("Select Sessions (Enter to toggle, Esc to exit)")
    ss.app.SetRoot(ss.list, true).SetFocus(ss.list).EnableMouse(true)

    ss.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyEsc {
            ss.updateGlobalSelectedSessions()
            ss.app.Stop()
            return nil
        }
        return event
    })

    if err := ss.app.Run(); err != nil {
        panic(err)
    }
}

func (ss *sessionSelector) updateGlobalSelectedSessions() {
    globals.Selected_Sessions = make([]string, 0, len(ss.selectedSessionIDs))
    for id := range ss.selectedSessionIDs {
        globals.Selected_Sessions = append(globals.Selected_Sessions, id)
    }
}

func SelectSessions() {
    selector := newSessionSelector()
    selector.run()
}

