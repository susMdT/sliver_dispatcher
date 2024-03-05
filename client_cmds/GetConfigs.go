package client_cmds

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

func GetConfigs(rpc rpcpb.SliverRPCClient) {
	fmt.Println("[Builds]")
	var builds *clientpb.ImplantBuilds
	builds, err := rpc.ImplantBuilds(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	var config *clientpb.ImplantConfig
	for _, config = range builds.GetConfigs() {
		fmt.Printf(fmt.Sprintf("ID: %-10s | GOOS: %-10s | ARCH: %-8s | C2: %-10s \n",
			strings.Split(config.ID, "-")[0],
			config.GOOS,
			config.GOARCH,
			config.C2))
	}

	fmt.Println("[Profiles]")
	var profiles *clientpb.ImplantProfiles
	profiles, err = rpc.ImplantProfiles(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	var profile *clientpb.ImplantProfile
	for _, profile = range profiles.GetProfiles() {
		fmt.Printf(fmt.Sprintf("%s\nID: %-10s | GOOS: %-10s | ARCH: %-8s | C2: %-10s \n",
			profile.Name,
			strings.Split(profile.Config.ID, "-")[0],
			profile.Config.GOOS,
			profile.Config.GOARCH,
			profile.Config.C2))
	}

}
