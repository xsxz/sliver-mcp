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
	Ifconfig = "session_run_ifconfig"
	Netstat  = "session_run_netstat"
)

var IfconfigTool = mcp.NewTool(
	Ifconfig,
	mcp.WithDescription("Get network interface configuration information about the remote system by session ID"),
	mcp.WithString("session_id", mcp.Required()),
)

func IfconfigHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Ifconfig(client, session)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run ifconfig command: %w", err)
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

////
////
////

var NetstatTool = mcp.NewTool(
	Netstat,
	mcp.WithDescription("Get network interface configuration information about the remote system by session ID"),
	mcp.WithString("session_id", mcp.Required()),
)

func NetstatHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Netstat(client, session)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat command: %w", err)
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
