package mcptools_sliver_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	ListListeners    = "listeners_list"
	KillAllListeners = "listeners_kill_all"

	StartHTTPListener = "listener_start_http"
	StartMTLSListener = "listener_start_mtls"
	StartWGListener   = "listener_start_wg"
	KillListener      = "listener_kill"
)

var ListListenersTool = mcp.NewTool(
	ListListeners,
	mcp.WithDescription("List all active listeners"),
)

func ListListenersHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	listeners, err := cmd.ListListeners(client)
	if err != nil {
		return nil, fmt.Errorf("failed to list listeners: %w", err)
	}

	if len(listeners) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("No active listeners found."),
			},
		}, nil
	}

	var sb strings.Builder
	sb.WriteString("Active listeners:\n")

	for _, job := range listeners {
		line := fmt.Sprintf(
			"ID: %d, Name: %s, Protocol: %s, Port: %d, Stage Profile: %s\n",
			job.ID, job.Name, job.Protocol, job.Port, job.ProfileName,
		)
		sb.WriteString(line)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(sb.String()),
		},
	}, nil
}

var KillAllListenersTool = mcp.NewTool(
	KillAllListeners,
	mcp.WithDescription("Kill all active listeners"),
)

func KillAllListenersHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	killedListeners, err := cmd.KillAllListeners(client)
	if err != nil {
		return nil, fmt.Errorf("failed to kill all listeners: %w", err)
	}

	var response string
	if len(killedListeners) == 0 {
		response = "No listeners were active."
	} else {
		response = "Killed the following listeners:\n" + strings.Join(killedListeners, "\n")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(response),
		},
	}, nil
}

////
////
////

var StartHTTPListenerTool = mcp.NewTool(
	StartHTTPListener,
	mcp.WithDescription("by session ID"),
	mcp.WithString("lhost", mcp.Required()),
	mcp.WithNumber("lport", mcp.Required()),
)

func StartHTTPListenerHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	lhost := request.GetString("lhost", "")
	lport := uint16(request.GetInt("lport", 4444))

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	msg, err := cmd.StartHTTPListener(client, lhost, lport)
	if err != nil {
		msg = err.Error()
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var StartMTLSListenerTool = mcp.NewTool(
	StartMTLSListener,
	mcp.WithDescription("Start mTLS listener"),
	mcp.WithString("lhost", mcp.Required()),
	mcp.WithNumber("lport", mcp.Required()),
)

func StartMTLSListenerHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	lhost := request.GetString("lhost", "")
	lport := uint16(request.GetInt("lport", 4444))

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	msg, err := cmd.StartMTLSListener(client, lhost, lport)
	if err != nil {
		msg = err.Error()
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var StartWGListenerTool = mcp.NewTool(
	StartWGListener,
	mcp.WithDescription("Start WireGuard listener"),
	mcp.WithNumber("lport", mcp.Required()),
	mcp.WithNumber("nport", mcp.Required()),
	mcp.WithNumber("keyExchangePort", mcp.Required()),
)

func StartWGListenerHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	lport := uint16(request.GetInt("lport", 0))
	nport := uint16(request.GetInt("nport", 0))
	keyPort := uint16(request.GetInt("keyExchangePort", 0))

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	msg, err := cmd.StartWGListener(client, lport, nport, keyPort)
	if err != nil {
		msg = err.Error()
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}

var KillListenerTool = mcp.NewTool(
	KillListener,
	mcp.WithDescription("Kill listener job by job ID"),
	mcp.WithNumber("job_id", mcp.Required()),
)

func KillListenerHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jobID := uint32(request.GetInt("job_id", 0))

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	msg, err := cmd.KillListener(client, jobID)
	if err != nil {
		msg = err.Error()
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(msg),
		},
	}, nil
}
