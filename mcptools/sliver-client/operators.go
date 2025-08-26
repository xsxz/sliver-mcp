package mcptools_sliver_client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
)

const (
	ListOperators = "operators_list"
)

var ListOperatorsTool = mcp.NewTool(
	ListOperators,
	mcp.WithDescription("List all remote operators, connected to the Sliver server"),
)

func ListOperatorsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	type OperatorInfo struct {
		Name   string `json:"name"`
		Online bool   `json:"online"`
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	operators, err := cmd.ListOperators(client)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve operators: %w", err)
	}

	var result []OperatorInfo
	for _, operator := range operators {
		result = append(result, OperatorInfo{
			Name:   operator.Name,
			Online: operator.Online,
		})
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operators: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonData)),
		},
	}, nil
}
