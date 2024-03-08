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

type ExtFile struct {
	OS   string
	Arch string
	Path string
}

type Arg struct {
	Name     string
	Type     string
	Optional bool
}
type ExtensionCfg struct {
	Command_Name string
	Entrypoint   string
	Files        []ExtFile
	Arguments    []Arg
}

var (
	DebugMode         bool
	ActiveSessions    []Interactive_Session
	Selected_Sessions []Interactive_Session
)
