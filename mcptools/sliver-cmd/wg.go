package mcptools_sliver_cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	ListWGSocks  = "session_wg_socks_list"
	StartWGSocks = "session_wg_socks_start"
	StopWGSocks  = "session_wg_socks_stop"

	ListWGPortFwd  = "session_wg_portfwd_list"
	StartWGPortFwd = "session_wg_portfwd_start"
	StopWGPortFwd  = "session_wg_portfwd_stop"
)

var ListWGSocksTool = mcp.NewTool(
	ListWGSocks,
	mcp.WithDescription("List active WireGuard SOCKS servers on a session"),
	mcp.WithString("session_id", mcp.Required()),
)

func ListWGSocksHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	type socksServer struct {
		ID        int32  `json:"id"`
		LocalAddr string `json:"local_address"`
	}

	var result []socksServer

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

	socksList, err := cmd.WGSocksList(client, session)
	if err != nil {
		return nil, fmt.Errorf("failed to list socks servers: %w", err)
	}

	if socksList.Servers != nil {
		for _, server := range socksList.Servers {
			result = append(result, socksServer{
				ID:        server.ID,
				LocalAddr: server.LocalAddr,
			})
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonBytes)),
		},
	}, nil
}

var StartWGSocksTool = mcp.NewTool(
	StartWGSocks,
	mcp.WithDescription("Start a WireGuard SOCKS server on a session"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("bind_port", mcp.Required()),
)

func StartWGSocksHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	bindPort := request.GetInt("bind_port", 0)
	if bindPort == 0 {
		return nil, fmt.Errorf("bind_port is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	socks, err := cmd.WGSocksStart(client, session, bindPort)
	if err != nil {
		return nil, fmt.Errorf("failed to start SOCKS server: %w", err)
	}

	if socks.Server == nil {
		return nil, fmt.Errorf("no SOCKS server returned")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Started SOCKS server on %s", socks.Server.LocalAddr)),
		},
	}, nil
}

var StopWGSocksTool = mcp.NewTool(
	StopWGSocks,
	mcp.WithDescription("Stop a WireGuard SOCKS server by ID on a session"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("socks_id", mcp.Required()),
)

func StopWGSocksHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	socksID := request.GetInt("socks_id", -1)
	if socksID == -1 {
		return nil, fmt.Errorf("socks_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	stopReq, err := cmd.WGSocksStop(client, session, socksID)
	if err != nil {
		return nil, fmt.Errorf("failed to stop SOCKS server: %w", err)
	}

	if stopReq.Server == nil {
		return nil, fmt.Errorf("no SOCKS server response received")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("Removed SOCKS listener rule on %s", stopReq.Server.LocalAddr)),
		},
	}, nil
}

var ListWGPortFwdTool = mcp.NewTool(
	ListWGPortFwd,
	mcp.WithDescription("List WireGuard port forwards on a session"),
	mcp.WithString("session_id", mcp.Required()),
)

func ListWGPortFwdHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	type portForward struct {
		ID         int32  `json:"id"`
		LocalAddr  string `json:"local_address"`
		RemoteAddr string `json:"remote_address"`
	}

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

	fwdList, err := cmd.WGPortFwdList(client, session)
	if err != nil {
		return nil, fmt.Errorf("failed to list port forwards: %w", err)
	}

	if fwdList.Response != nil && fwdList.Response.Err != "" {
		return nil, fmt.Errorf("error: %s", fwdList.Response.Err)
	}

	var forwards []portForward
	if fwdList.Forwarders != nil {
		for _, fwd := range fwdList.Forwarders {
			forwards = append(forwards, portForward{
				ID:         fwd.ID,
				LocalAddr:  fwd.LocalAddr,
				RemoteAddr: fwd.RemoteAddr,
			})
		}
	}

	jsonBytes, err := json.Marshal(forwards)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonBytes)),
		},
	}, nil
}

var StartWGPortFwdTool = mcp.NewTool(
	StartWGPortFwd,
	mcp.WithDescription("Add a WireGuard port forward on a session"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("local_port", mcp.Required()),
	mcp.WithString("remote_addr", mcp.Required()),
)

func StartWGPortFwdHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	localPort := request.GetInt("local_port", 0)
	if localPort == 0 {
		return nil, fmt.Errorf("local_port is required")
	}

	remoteAddr := request.GetString("remote_addr", "")
	if remoteAddr == "" {
		return nil, fmt.Errorf("remote_addr is required")
	}

	_, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote target %s: %w", remoteAddr, err)
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	portfwdAdd, err := cmd.WGPortFwdAdd(client, session, localPort, remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to add port forward: %w", err)
	}

	if portfwdAdd.Response != nil && portfwdAdd.Response.Err != "" {
		return nil, fmt.Errorf("error: %s", portfwdAdd.Response.Err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf(
				"Port forwarding from %s to %s",
				portfwdAdd.Forwarder.LocalAddr,
				remoteAddr,
			)),
		},
	}, nil
}

var StopWGPortFwdTool = mcp.NewTool(
	StopWGPortFwd,
	mcp.WithDescription("Remove a WireGuard port forward by ID on a session"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithNumber("forward_id", mcp.Required()),
)

func StopWGPortFwdHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	forwardID := request.GetInt("forward_id", -1)
	if forwardID == -1 {
		return nil, fmt.Errorf("forward_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	stopReq, err := cmd.WGPortFwdRm(client, session, forwardID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove port forward: %w", err)
	}

	if stopReq.Response != nil && stopReq.Response.Err != "" {
		return nil, fmt.Errorf("error: %s", stopReq.Response.Err)
	}

	if stopReq.Forwarder == nil {
		return nil, fmt.Errorf("no forwarder info returned")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf(
				"Removed port forwarding rule %s -> %s",
				stopReq.Forwarder.LocalAddr,
				stopReq.Forwarder.RemoteAddr,
			)),
		},
	}, nil
}
