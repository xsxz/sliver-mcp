package sliver_cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
	"time"

	gen "github.com/bishopfox/sliver/client/command/generate"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/util"

	c "github.com/xsxz/sliver-mcp/sliver/client"
	"github.com/xsxz/sliver-mcp/utils"
)

func ListImplantBuilds(client *c.SliverClient) (*clientpb.ImplantBuilds, error) {
	builds, err := client.RPC.ImplantBuilds(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, err
	}

	if len(builds.Configs) == 0 {
		return nil, fmt.Errorf("no implant builds found")
	}

	return builds, nil
}

func RmImplantBuild(client *c.SliverClient, implantName string) (string, error) {
	builds, err := ListImplantBuilds(client)
	if err != nil {
		return "", err
	}

	if _, exists := builds.Configs[implantName]; !exists {
		return "", fmt.Errorf("implant build '%s' does not exist", implantName)
	}

	_, err = client.RPC.DeleteImplantBuild(context.Background(), &clientpb.DeleteReq{
		Name: implantName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to delete implant: %v", err)
	}

	return fmt.Sprintf("Implant build '%s' was deleted", implantName), nil
}

////
////
////

func ListImplantProfiles(client *c.SliverClient) (*clientpb.ImplantProfiles, error) {
	pbProfiles, err := client.RPC.ImplantProfiles(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, err
	}

	return pbProfiles, nil
}

// for demo purposes, generated profile to be used with backdoor command only
func GenerateWinShellcodeImplantProfile(client *c.SliverClient, name string, c2type string, connStrings string) (string, error) {
	validC2Types := []string{"mtls", "wg", "http", "named-pipe", "tcp-pivot"}

	var tunIP net.IP
	type WGConfig struct {
		WGPeerTunIP       string
		WGKeyExchangePort uint32
		WGTcpCommsPort    uint32
	}

	WGOptions := WGConfig{
		WGPeerTunIP:       "",
		WGKeyExchangePort: 0,
		WGTcpCommsPort:    0,
	}

	if name != "" {
		name = strings.ToLower(name)

		if err := util.AllowedName(name); err != nil {
			return "", err
		}
	}

	if !slices.ContainsFunc(validC2Types, func(t string) bool {
		return strings.EqualFold(c2type, t)
	}) {
		return "", errors.New("invalid C2 type: should be mtls, wg, http, named-pipe or tcp-pivot")
	}

	c2s := []*clientpb.ImplantC2{}

	if strings.EqualFold(c2type, "mtls") {
		mtlsC2, err := gen.ParseMTLSc2(connStrings)
		if err != nil {
			return "", err
		}
		c2s = append(c2s, mtlsC2...)
	}

	if strings.EqualFold(c2type, "wg") {
		mtlsC2, err := gen.ParseWGc2(connStrings)
		if err != nil {
			return "", err
		}
		c2s = append(c2s, mtlsC2...)

		uniqueWGIP, err := client.RPC.GenerateUniqueIP(context.Background(), &commonpb.Empty{})
		tunIP = net.ParseIP(uniqueWGIP.IP)
		if err != nil {
			return "", fmt.Errorf("failed to generate unique ip for wg peer tun interface")
		}

	}

	if strings.EqualFold(c2type, "http") {
		mtlsC2, err := gen.ParseHTTPc2(connStrings)
		if err != nil {
			return "", err
		}
		c2s = append(c2s, mtlsC2...)
	}

	if strings.EqualFold(c2type, "named-pipe") {
		mtlsC2, err := gen.ParseNamedPipec2(connStrings)
		if err != nil {
			return "", err
		}
		c2s = append(c2s, mtlsC2...)
	}

	if strings.EqualFold(c2type, "tcp-pivot") {
		mtlsC2, err := gen.ParseTCPPivotc2(connStrings)
		if err != nil {
			return "", err
		}
		c2s = append(c2s, mtlsC2...)
	}

	config := &clientpb.ImplantConfig{
		GOOS:             "windows",
		GOARCH:           "amd64",
		Name:             name,
		Debug:            false,
		Evasion:          true,
		ObfuscateSymbols: true,
		C2:               c2s,
		CanaryDomains:    nil,
		TemplateName:     "",

		WGPeerTunIP:       tunIP.String(),
		WGKeyExchangePort: WGOptions.WGKeyExchangePort,
		WGTcpCommsPort:    WGOptions.WGTcpCommsPort,

		ConnectionStrategy:  "s", //connectionStrategy,
		ReconnectInterval:   int64(utils.DefaultTimeout) * int64(time.Second),
		PollTimeout:         int64(utils.DefaultTimeout*6) * int64(time.Second), // 360, default
		MaxConnectionErrors: uint32(1000),

		LimitDomainJoined: false,
		LimitHostname:     "",
		LimitUsername:     "",
		LimitDatetime:     "",
		LimitFileExists:   "",
		LimitLocale:       "",

		Format:      clientpb.OutputFormat_SHELLCODE,
		IsSharedLib: false,
		IsService:   false,
		IsShellcode: true,

		RunAtLoad: false,

		DebugFile: "",
	}

	profile := &clientpb.ImplantProfile{
		Name:   name,
		Config: config,
	}
	resp, err := client.RPC.SaveImplantProfile(context.Background(), profile)
	if err != nil {
		return "", err
	}

	//fmt.Println(config.C2)

	// arch: default amd64 (hardcoded)
	// os: windows (hardcoded)
	// format: shellcode (hardcoded)
	// reconnect 60 secs (hardcoded)

	// name
	// disable sgn: bool
	// evasion: bool

	// http: connection strings
	// mtls: connection strings
	// wg: connection strings, key exchange port, c2 comms port
	// tcp-pivot connection strings
	// named pipe: connection strings

	return fmt.Sprintf("Saved new implant profile (beacon) %s", resp.Name), nil
}

func RmImplantProfile(client *c.SliverClient, profileName string) (string, error) {
	profiles, err := ListImplantProfiles(client)
	if err != nil {
		return "", err
	}

	if !slices.ContainsFunc(profiles.Profiles, func(p *clientpb.ImplantProfile) bool {
		return p.Name == profileName
	}) {
		return "", fmt.Errorf("profile '%s' does not exist", profileName)
	}

	_, err = client.RPC.DeleteImplantProfile(context.Background(), &clientpb.DeleteReq{
		Name: profileName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to delete profile: %w", err)
	}

	return fmt.Sprintf("removed implant profile: %s", profileName), nil
}
