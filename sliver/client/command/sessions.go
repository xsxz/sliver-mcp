package sliver_cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
	"github.com/xsxz/sliver-mcp/utils"
)

func GetSessionImplantConfig(client *c.SliverClient, sessionID string) (*clientpb.ImplantConfig, error) {
	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)

	}

	c2s := []*clientpb.ImplantC2{}
	c2s = append(c2s, &clientpb.ImplantC2{
		URL:      session.GetActiveC2(),
		Priority: uint32(0),
	})
	config := &clientpb.ImplantConfig{
		ID:      session.ID,
		Name:    session.GetName(),
		GOOS:    session.GetOS(),
		GOARCH:  session.GetArch(),
		Debug:   true,
		Evasion: session.GetEvasion(),

		MaxConnectionErrors: uint32(1000),
		ReconnectInterval:   int64(60),
		Format:              clientpb.OutputFormat_SHELLCODE,
		IsSharedLib:         true,
		C2:                  c2s,
	}

	return config, nil
}

func KillSession(client *c.SliverClient, sessionID string, force bool) (string, error) {
	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve session: %w", err)
	}

	_, err = client.RPC.Kill(context.Background(), &sliverpb.KillReq{
		Request: &commonpb.Request{
			SessionID: session.ID,
			Timeout:   int64(utils.DefaultTimeout),
		},
		Force: force,
	})
	if err != nil {
		return "", fmt.Errorf("failed to kill session: %w", err)
	}

	return fmt.Sprintf("Session %s was killed", sessionID), nil
}

// CloseInteractiveSession - Close an interactive session but do not kill the remote process
func CloseInteractiveSession(client *c.SliverClient, sessionID string) (string, error) {
	s, err := client.GetSessionByID(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve session: %w", err)
	}

	_, err = client.RPC.CloseSession(context.Background(), &sliverpb.CloseSession{
		Request: &commonpb.Request{
			SessionID: s.ID,
			Timeout:   int64(utils.DefaultTimeout),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to close session: %w", err)
	}

	return fmt.Sprintf("Session %s was closed", sessionID), nil
}

func PruneSessions(client *c.SliverClient) (string, error) {
	sessions, err := client.RPC.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		return "", err
	}

	if len(sessions.GetSessions()) == 0 {
		return "No sessions to prune", nil
	}

	var prunedSessions []string

	for _, session := range sessions.GetSessions() {
		if session.IsDead {
			_, err = KillSession(client, session.ID, true)
			if err == nil {
				prunedSessions = append(prunedSessions, session.ID)
			}
			// currently errors are skipped
		}
	}

	if len(prunedSessions) == 0 {
		return "No sessions were pruned", nil
	}

	return fmt.Sprintf("Pruned dead sessions:\n%s", strings.Join(prunedSessions, " ")), nil
}
