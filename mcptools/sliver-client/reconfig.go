package mcptools_sliver_client

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	RenameImplant = "implant_rename"
)

var RenameImplantTool = mcp.NewTool(
	RenameImplant,
	mcp.WithDescription("Rename target session or beacon"),
	mcp.WithString("session_id"),
	mcp.WithString("beacon_id"),
	mcp.WithString("implant_name", mcp.Required()),
)

func RenameImplantHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	beaconID := request.GetString("beacon_id", "")
	implantName := request.GetString("implant_name", "")
	if implantName == "" {
		return nil, fmt.Errorf("implant_name is required")
	}

	if (sessionID == "" && beaconID == "") || (sessionID != "" && beaconID != "") {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Error: must provide either session ID or beacon ID, not both"),
			},
		}, nil
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	result, err := cmd.RenameImplant(client, sessionID, beaconID, implantName)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil

}
