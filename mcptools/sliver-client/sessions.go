package mcptools_sliver_client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	ListSessions    = "sessions_list"
	KillAllSessions = "sessions_kill_all"
	PruneSessions   = "sessions_prune"

	GetSessionInfo          = "session_get_info"
	CloseInteractiveSession = "session_close"
	KillSession             = "session_kill"
)

var ListSessionsTool = mcp.NewTool(
	ListSessions,
	mcp.WithDescription("List all available Sliver sessions"),
)

func ListSessionsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var activeSessions []string

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	sessions, err := client.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, session := range sessions.Sessions {
		activeSessions = append(activeSessions, session.ID)
	}

	result := fmt.Sprintf("Active sessions: %s", strings.Join(activeSessions, " "))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

var GetSessionInfoTool = mcp.NewTool(
	GetSessionInfo,
	mcp.WithDescription("Get Sliver session info by ID"),
	mcp.WithString("session_id", mcp.Required()),
)

func GetSessionInfoHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.GetSessionInfo(client, sessionID)
	}

	sessionInfo, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	sessionMap, ok := sessionInfo.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected sessionInfo type: %T", sessionInfo)
	}

	jsonBytes, err := json.Marshal(sessionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session info to JSON: %w", err)
	}

	text := string(jsonBytes)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(text),
		},
	}, nil
}

////
////
////

var CloseInteractiveSessionTool = mcp.NewTool(
	CloseInteractiveSession,
	mcp.WithDescription("Close an interactive session but do not kill the remote process. by session ID"),
	mcp.WithString("session_id", mcp.Required()),
)

func CloseInteractiveSessionHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	result, err := cmd.CloseInteractiveSession(client, sessionID)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

////
////
////

var PruneSessionsTool = mcp.NewTool(
	PruneSessions,
	mcp.WithDescription("Prune dead sessions"),
)

func PruneSessionsToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	result, err := cmd.PruneSessions(client)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

var KillSessionTool = mcp.NewTool(
	KillSession,
	mcp.WithDescription("Kill beacon by ID"),
	mcp.WithString("beacon_id", mcp.Required()),
)

func KillSessionHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	result, err := cmd.KillSession(client, sessionID, true)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

var KillAllSessionsTool = mcp.NewTool(
	KillAllSessions,
	mcp.WithDescription("Kill all interactive sessions"),
)

func KillAllSessionsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var killedSessions []string

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	sessions, err := client.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, session := range sessions.Sessions {
		msg, err := cmd.KillSession(client, session.ID, true)
		if err != nil {
			fmt.Printf("Failed to kill session %s: %s\n", session.ID, err)
			continue
		}
		killedSessions = append(killedSessions, msg)
	}

	result := "No sessions were killed."
	if len(killedSessions) > 0 {
		result = fmt.Sprintf("Killed sessions: %s", strings.Join(killedSessions, " "))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}
