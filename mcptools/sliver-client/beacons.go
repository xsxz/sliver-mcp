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
	ListBeacons    = "beacons_list"
	KillAllBeacons = "beacons_kill_all"
	PruneBeacons   = "beacons_prune"

	GetBeaconInfo                 = "beacon_get_info"
	StartBeaconInteractiveSession = "beacon_start_interactive_session"
	KillBeacon                    = "beacon_kill"
)

var ListBeaconsTool = mcp.NewTool(
	ListBeacons,
	mcp.WithDescription("List all available Sliver beacons"),
)

func ListBeaconsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var activeBeacons []string

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	beacons, err := client.ListBeacons()
	if err != nil {
		return nil, fmt.Errorf("failed to list beacons: %w", err)
	}

	for _, beacon := range beacons.Beacons {
		activeBeacons = append(activeBeacons, beacon.ID)
	}

	var result string
	if len(activeBeacons) == 0 {
		result = "No active beacons."
	} else {
		result = fmt.Sprintf("Active beacons: %s", strings.Join(activeBeacons, " "))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

var GetBeaconInfoTool = mcp.NewTool(
	GetBeaconInfo,
	mcp.WithDescription("Get Sliver beacon info by ID"),
	mcp.WithString("beacon_id", mcp.Required()),
)

func GetBeaconInfoHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	beaconID := request.GetString("beacon_id", "")
	if beaconID == "" {
		return nil, fmt.Errorf("beacon_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(beacon *clientpb.Beacon) (interface{}, error) {
		return cmd.GetBeaconInfo(client, beaconID)
	}

	result, err := client.RunInBeacon(beaconID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	jsonBytes, err := json.Marshal(resultMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal beacon info to JSON: %w", err)
	}

	text := string(jsonBytes)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(text),
		},
	}, nil
}

var StartBeaconInteractiveSessionTool = mcp.NewTool(
	StartBeaconInteractiveSession,
	mcp.WithDescription("Start interactive session from beacon. by beacon ID"),
	mcp.WithString("beacon_id", mcp.Required()),
)

func StartBeaconInteractiveSessionHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	beaconID := request.GetString("beacon_id", "")
	if beaconID == "" {
		return nil, fmt.Errorf("beacon_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	response, err := cmd.StartBeaconInteractiveSession(client, beaconID)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal beacon info to JSON: %w", err)
	}

	text := string(jsonBytes)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(text),
		},
	}, nil
}

var KillBeaconTool = mcp.NewTool(
	KillBeacon,
	mcp.WithDescription("Kill beacon by ID"),
	mcp.WithString("beacon_id", mcp.Required()),
)

func KillBeaconHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	beaconID := request.GetString("beacon_id", "")
	if beaconID == "" {
		return nil, fmt.Errorf("beacon_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	result, err := cmd.KillBeacon(client, beaconID, true)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

var KillAllBeaconsTool = mcp.NewTool(
	KillAllBeacons,
	mcp.WithDescription("Kill all beacons"),
)

func KillAllBeaconsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var killedBeacons []string

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	beacons, err := client.ListBeacons()
	if err != nil {
		return nil, fmt.Errorf("failed to list beacons: %w", err)
	}

	for _, beacon := range beacons.Beacons {
		_, err := cmd.KillBeacon(client, beacon.ID, true)
		if err != nil {
			fmt.Printf("Failed to kill beacon %s: %s\n", beacon.ID, err)
			continue
		}
		killedBeacons = append(killedBeacons, beacon.ID)
	}

	result := fmt.Sprintf("Killed beacons: %s", strings.Join(killedBeacons, " "))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

var PruneBeaconsTool = mcp.NewTool(
	PruneBeacons,
	mcp.WithDescription("Prune stale beacons"),
	mcp.WithString("duration"),
)

func PruneBeaconsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	duration := request.GetString("duration", "1h")

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	pruneResult, err := cmd.PruneBeacons(client, duration)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(pruneResult),
		},
	}, nil
}
