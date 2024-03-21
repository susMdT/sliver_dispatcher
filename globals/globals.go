package globals

import (
	"sync"

	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

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

// TODO: when implementing selective dispatch, pass these codes instead of an OS string
const (
	DISPATCH_WIN    = 1
	DISPATCH_LIN    = 2
	DISPATCH_CUSTOM = 3
)

var (
	L_rpc             sync.Mutex
	DebugMode         bool
	ActiveSessions    []Interactive_Session
	Selected_Sessions []string
	Rpc               rpcpb.SliverRPCClient
	DispatchType      int
)
