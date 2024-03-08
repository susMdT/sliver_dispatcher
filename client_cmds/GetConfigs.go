package client_cmds

import (
	"context"
	"sliver-dispatch/utils"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

func GetConfigs(rpc rpcpb.SliverRPCClient) {
	utils.Iprint("[Builds]")
	var builds *clientpb.ImplantBuilds
	builds, err := rpc.ImplantBuilds(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error listing implant builds: %s", err)
	}
	var config *clientpb.ImplantConfig
	for _, config = range builds.GetConfigs() {
		utils.Iprint("ID: %-10s | GOOS: %-10s | ARCH: %-8s | C2: %-10s",
			strings.Split(config.ID, "-")[0],
			config.GOOS,
			config.GOARCH,
			config.C2)
	}

	utils.Iprint("[Profiles]")
	var profiles *clientpb.ImplantProfiles
	profiles, err = rpc.ImplantProfiles(context.Background(), &commonpb.Empty{})
	if err != nil {
		utils.Eprint("Error listing implant profiles: %s", err)
	}
	var profile *clientpb.ImplantProfile
	for _, profile = range profiles.GetProfiles() {
		utils.Iprint("%s\nID: %-10s | GOOS: %-10s | ARCH: %-8s | C2: %-10s",
			profile.Name,
			strings.Split(profile.Config.ID, "-")[0],
			profile.Config.GOOS,
			profile.Config.GOARCH,
			profile.Config.C2)
	}

}
