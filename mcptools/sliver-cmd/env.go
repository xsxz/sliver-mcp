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
	SetEnv   = "session_run_setenv"
	GetEnv   = "session_run_getenv"
	UnsetEnv = "session_run_unsetenv"
)

var SetEnvTool = mcp.NewTool(
	SetEnv,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("var_name", mcp.Required()),
	mcp.WithString("var_value", mcp.Required()),
)

func SetEnvHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	varName := request.GetString("var_name", "")
	varValue := request.GetString("var_value", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.SetEnv(client, session, varName, varValue)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
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

var GetEnvTool = mcp.NewTool(
	GetEnv,
	mcp.WithDescription("get env of process (variable can be empty to get all the vars at once; specify to get only key=value), by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("var_name", mcp.Required()),
)

func GetEnvHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	varName := request.GetString("var_name", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.GetEnv(client, session, varName)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
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

var UnsetEnvTool = mcp.NewTool(
	UnsetEnv,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
)

func UnsetEnvHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	varName := request.GetString("var_name", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.UnsetEnv(client, session, varName)
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
