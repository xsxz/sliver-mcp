package mcptools_sliver_cmd

import (
	"context"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	Backdoor = "session_run_backdoor"
)

var BackdoorTool = mcp.NewTool(
	Backdoor,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("profile_name", mcp.Required()),
	mcp.WithString("remote_file_path", mcp.Required()),
)

func BackdoorHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	remote_file_path := request.GetString("remote_file_path", "")
	profile_name := request.GetString("profile_name", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Backdoor(client, sessionID, profile_name, remote_file_path)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result.(string)),
		},
	}, nil
}
