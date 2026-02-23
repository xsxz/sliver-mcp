package mcptools_sliver_cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	ListPivots    = "session_pivots_list"
	StopAllPivots = "session_pivots_stop_all"

	StartTcpPivot       = "session_pivot_start_tcp"
	StartNamedPipePivot = "session_pivot_start_namedpipe"
	StopPivot           = "session_pivot_stop"
)

var StartTcpPivotTool = mcp.NewTool(
	StartTcpPivot,
	mcp.WithDescription("Start TCP pivot listener"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("bind_address", mcp.Required()),
	mcp.WithNumber("lport", mcp.Required()),
)

func StartTcpPivotHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	bindAddress := request.GetString("bind_address", "")
	lport := uint16(request.GetInt("lport", 9898))

	if sessionID == "" || bindAddress == "" || lport == 0 {
		return nil, fmt.Errorf("session_id, bind_address and lport are required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	msg, err := cmd.StartTcpPivot(client, session, bindAddress, lport)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var StartNamedPipePivotTool = mcp.NewTool(
	StartNamedPipePivot,
	mcp.WithDescription("Start named pipe pivot listener"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("bind_address", mcp.Required()),
	mcp.WithBoolean("allow_all"),
)

func StartNamedPipePivotHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	bindAddress := request.GetString("bind_address", "")
	allowAll := request.GetBool("allow_all", false)

	if sessionID == "" || bindAddress == "" {
		return nil, fmt.Errorf("session_id and bind_address are required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	msg, err := cmd.StartNamedPipePivot(client, session, bindAddress, allowAll)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var StopPivotTool = mcp.NewTool(
	StopPivot,
	mcp.WithDescription("Stop pivot listener by ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("id", mcp.Required()),
)

func StopPivotHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	id := uint32(request.GetInt("id", 0))

	if sessionID == "" || id == 0 {
		return nil, fmt.Errorf("session_id and id are required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	msg, err := cmd.StopPivot(client, session, id)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var StopAllPivotsTool = mcp.NewTool(
	StopAllPivots,
	mcp.WithDescription("Stop all pivot listeners on a session"),
	mcp.WithString("session_id", mcp.Required()),
)

func StopAllPivotsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")

	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	msg, err := cmd.StopAllPivots(client, session)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var ListPivotsTool = mcp.NewTool(
	ListPivots,
	mcp.WithDescription("Stop all pivot listeners on a session"),
	mcp.WithString("session_id", mcp.Required()),
)

func ListPivotsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	type pivotDetail struct {
		ID          uint32 `json:"id"`
		Protocol    string `json:"protocol"`
		BindAddress string `json:"bind_address"`
		NumOfPivots int    `json:"num_of_pivots"`
	}
	var details []pivotDetail

	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	pivotListeners, err := cmd.ListPivots(client, session)
	if err != nil {
		return nil, err
	}

	for _, listener := range pivotListeners.Listeners {
		details = append(details, pivotDetail{
			ID:          listener.ID,
			Protocol:    listener.Type.String(),
			BindAddress: listener.BindAddress,
			NumOfPivots: len(listener.Pivots),
		})
	}
	jsonBytes, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonBytes)),
		},
	}, nil
}
