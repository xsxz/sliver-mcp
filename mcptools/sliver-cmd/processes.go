package mcptools_sliver_cmd

import (
	"context"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/golang/protobuf/proto"
	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
	"github.com/xsxz/sliver-mcp/utils"
)

const (
	Ps        = "session_run_ps"
	Terminate = "session_run_terminate_proc"
)

var PsTool = mcp.NewTool(
	Ps,
	mcp.WithDescription("Get list of running processes on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
)

func PsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Ps(client, session)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run ps command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var TerminateTool = mcp.NewTool(
	Terminate,
	mcp.WithDescription("Forcefully terminate process on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("pid", mcp.Required()),
)

func TerminateHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pid := int32(request.GetInt("pid", 0))
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Terminate(client, session, pid)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run terminate command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}
