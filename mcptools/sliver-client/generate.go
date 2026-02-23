package mcptools_sliver_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	ListImplantBuilds = "implants_builds_list_all"
	RmImplantBuild    = "implant_build_rm"

	ListImplantProfiles         = "implants_profiles_list_all"
	GenerateWinShellcodeProfile = "implant_profile_generate_winshellcode"
	RmImplantProfile            = "implant_profile_rm"
)

var ListImplantBuildsTool = mcp.NewTool(
	ListImplantBuilds,
	mcp.WithDescription("List all implant builds"),
)

func ListImplantBuildsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	builds, err := cmd.ListImplantBuilds(client)
	if err != nil {
		return nil, fmt.Errorf("failed to list implant builds: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("Implant Builds:")
	for name, cfg := range builds.Configs {
		line := fmt.Sprintf(
			"Name: %s, Beacon: %t, Template: %s, OS/Arch: %s/%s, Format: %s, C2: %s, Debug: %t",
			name, cfg.IsBeacon, cfg.TemplateName, cfg.GOOS, cfg.GOARCH, cfg.Format, cfg.C2, cfg.Debug,
		)
		sb.WriteString(line)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(sb.String()),
		},
	}, nil
}

var RmImplantBuildTool = mcp.NewTool(
	RmImplantBuild,
	mcp.WithDescription("Remove an implant build by name"),
	mcp.WithString("implant_build_name", mcp.Required()),
)

func RmImplantBuildHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	implantName := request.GetString("implant_name", "")
	if implantName == "" {
		return nil, fmt.Errorf("implant_name is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	_, err = cmd.RmImplantBuild(client, implantName)
	if err != nil {
		return nil, fmt.Errorf("failed to remove implant build: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Successfully removed implant build: %s", implantName)),
		},
	}, nil
}

///
///
///

var ListImplantProfilesTool = mcp.NewTool(
	ListImplantProfiles,
	mcp.WithDescription("List all implant profiles"),
)

func ListImplantProfilesHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	resp, err := cmd.ListImplantProfiles(client)
	if err != nil {
		return nil, fmt.Errorf("failed to list implant profiles: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("Implant Profiles: ")

	for _, profile := range resp.Profiles {
		cfg := profile.Config
		if cfg == nil {
			continue
		}

		implantType := "session"
		if cfg.IsBeacon {
			implantType = "beacon"
		}

		obfuscation := "disabled"
		if cfg.ObfuscateSymbols {
			obfuscation = "enabled"
		}

		c2List := []string{}
		for i, c2 := range cfg.C2 {
			c2List = append(c2List, fmt.Sprintf("[%d] %s", i+1, c2.URL))
		}
		c2Str := strings.Join(c2List, ", ")

		line := fmt.Sprintf(
			"Name: %s, Type: %s,  Template: %s, OS/Arch: %s/%s, Format: %s Debug: %t, Obfuscation: %s, C2: %s",
			profile.Name,
			implantType,
			cfg.TemplateName,
			cfg.GOOS,
			cfg.GOARCH,
			cfg.Format.String(),
			cfg.Debug,
			obfuscation,
			c2Str,
		)

		sb.WriteString(line)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(sb.String()),
		},
	}, nil
}

var RmImplantProfileTool = mcp.NewTool(
	RmImplantProfile,
	mcp.WithDescription("Remove an implant profile by name"),
	mcp.WithString("implant_profile_name", mcp.Required()),
)

func RmImplantProfileHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	implantProfileName := request.GetString("implant_profile_name", "")
	if implantProfileName == "" {
		return nil, fmt.Errorf("implant_name is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	response, err := cmd.RmImplantProfile(client, implantProfileName)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(response),
		},
	}, nil
}

var GenerateWinShellcodeProfileTool = mcp.NewTool(
	GenerateWinShellcodeProfile,
	mcp.WithDescription("Generate a Windows shellcode implant profile"),
	mcp.WithString("profile_name", mcp.Required()),
	mcp.WithString("c2_type", mcp.Required()),
	mcp.WithString("conn_strings", mcp.Required()),
)

func GenerateWinShellcodeProfileHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	profile_name := request.GetString("profile_name", "")
	c2type := request.GetString("c2_type", "")
	connStrings := request.GetString("conn_strings", "")

	if profile_name == "" || c2type == "" || connStrings == "" {
		return nil, fmt.Errorf("name, c2_type, and conn_strings are required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	response, err := cmd.GenerateWinShellcodeImplantProfile(client, profile_name, c2type, connStrings)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shellcode implant profile: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(response),
		},
	}, nil
}
