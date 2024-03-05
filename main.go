package main

import (
	"sliver-dispatch/tui"
	"sliver-dispatch/utils"

	"flag"
	"log"

	"github.com/bishopfox/sliver/client/assets"
	"github.com/bishopfox/sliver/client/transport"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "/root/.sliver-client/configs/default.cfg", "path to sliver client config file")
	flag.Parse()

	// load the client configuration from the filesystem
	config, err := assets.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	// connect to the server
	var rpc rpcpb.SliverRPCClient
	rpc, ln, err := transport.MTLSConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	utils.UpdateSessions(rpc)

	tui.Main(rpc)
}
