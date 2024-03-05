package globals

type Interactive_Session struct {
	Host     string
	Address  string
	OS       string
	ID       string
	PID      int
	Process  string
	Username string
}

var (
	ActiveSessions    []Interactive_Session
	Selected_Sessions []Interactive_Session
)
