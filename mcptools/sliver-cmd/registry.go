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
	RegReadValue      = "session_registry_read_value"
	RegReadKeysValues = "session_registry_read_key_values"
	RegCreateKey      = "session_registry_create_key"
	RegWriteValue     = "session_registry_write_value"
	RegDeleteKey      = "session_registry_delete_key"
)

var RegReadKeysValuesTool = mcp.NewTool(
	RegReadKeysValues,
	mcp.WithDescription("Read all key values from registry hive on the remote Windows system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("hive", mcp.Required()),
	mcp.WithString("reg_path", mcp.Required()),
)

func RegReadKeysValuesHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hive := request.GetString("hive", "")
	regPath := request.GetString("reg_path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.RegReadKeysValues(client, session, hive, regPath)
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

var RegReadValueTool = mcp.NewTool(
	RegReadValue,
	mcp.WithDescription("Read one key value from registry hive on the remote Windows system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("hive", mcp.Required()),
	mcp.WithString("reg_path", mcp.Required()),
)

func RegReadValueHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hive := request.GetString("hive", "")
	regPath := request.GetString("reg_path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.RegReadValue(client, session, hive, regPath)
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

var RegCreateKeyTool = mcp.NewTool(
	RegCreateKey,
	mcp.WithDescription("Get list of running processes on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("hive", mcp.Required()),
	mcp.WithString("reg_path", mcp.Required()),
)

func RegCreateKeyHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hive := request.GetString("hive", "")
	regPath := request.GetString("reg_path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.RegCreateKey(client, session, hive, regPath)
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

var RegWriteValueTool = mcp.NewTool(
	RegWriteValue,
	mcp.WithDescription("Get list of running processes on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("hive", mcp.Required()),
	mcp.WithString("reg_path", mcp.Required()),
	mcp.WithString("flag_type", mcp.Required()),
	mcp.WithString("value", mcp.Required()),
)

func RegWriteValueHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hive := request.GetString("hive", "")
	regPath := request.GetString("reg_path", "")
	flagType := request.GetString("flag_type", "")
	value := request.GetString("value", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.RegWriteValue(client, session, hive, regPath, flagType, value)
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

var RegDeleteKeyTool = mcp.NewTool(
	RegDeleteKey,
	mcp.WithDescription("Get list of running processes on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("hive", mcp.Required()),
	mcp.WithString("reg_path", mcp.Required()),
)

func RegDeleteKeyHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hive := request.GetString("hive", "")
	regPath := request.GetString("reg_path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.RegDeleteKey(client, session, hive, regPath)
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
