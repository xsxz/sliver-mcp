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
	GetPrivs    = "session_run_getprivs"
	GetSystem   = "session_run_getsystem"
	Impersonate = "session_run_impersonate"
	MakeToken   = "session_run_maketoken"
	Rev2Self    = "session_run_rev2self"
	RunAs       = "session_run_runas"
)

var GetPrivsTool = mcp.NewTool(
	GetPrivs,
	mcp.WithDescription("Get list of running processes on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
)

func GetPrivsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.GetPrivs(client, session)
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

var GetSystemTool = mcp.NewTool(
	GetSystem,
	mcp.WithDescription("Forcefully terminate process on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("process_name_inject_to"),
)

func GetSystemHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	processNameInjectTo := request.GetString("process_name_inject_to", "spoolsv.exe")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.GetSystem(client, session, processNameInjectTo)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run terminate command: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result.(string)),
		},
	}, nil
}

var ImpersonateTool = mcp.NewTool(
	Impersonate,
	mcp.WithDescription("Forcefully terminate process on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("username", mcp.Required()),
)

func ImpersonateHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	username := request.GetString("username", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Impersonate(client, session, username)
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

var MakeTokenTool = mcp.NewTool(
	MakeToken,
	mcp.WithDescription("Forcefully terminate process on the remote system"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("username", mcp.Required()),
	mcp.WithString("password", mcp.Required()),
	mcp.WithString("domain"),
	mcp.WithString("logon_type", mcp.Required()),
)

func MakeTokenHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	username := request.GetString("username", "")
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	password := request.GetString("password", "")
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	logonType := request.GetString("logon_type", "")
	if logonType == "" {
		return nil, fmt.Errorf("logon type is required")
	}

	domain := request.GetString("domain", "")

	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.MakeToken(client, session, username, password, domain, logonType)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run terminate command: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result.(string)),
		},
	}, nil
}

var Rev2SelfTool = mcp.NewTool(
	Rev2Self,
	mcp.WithDescription("Drop all impersonated tokens"),
	mcp.WithString("session_id", mcp.Required()),
)

func Rev2SelfHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Rev2Self(client, session)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run terminate command: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result.(string)),
		},
	}, nil
}

var RunAsTool = mcp.NewTool(
	RunAs,
	mcp.WithDescription("Drop all impersonated tokens"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("username", mcp.Required()),
	mcp.WithString("password", mcp.Required()),
	mcp.WithString("domain"),
	mcp.WithString("args"),
	mcp.WithString("netOnly"),
	mcp.WithString("hideWindow"),
)

func RunAsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	username := request.GetString("username", "")
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	password := request.GetString("password", "")
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	processName := request.GetString("process_path", "")
	if sessionID == "" {
		return nil, fmt.Errorf("process_path is required")
	}

	domain := request.GetString("domain", "")
	args := request.GetString("args", "")
	netOnly := request.GetBool("netOnly", false)
	hideWindow := request.GetBool("hideWindow", false)

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.RunAs(client, session, username, password, domain, processName, args, netOnly, hideWindow)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run terminate command: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result.(string)),
		},
	}, nil
}
