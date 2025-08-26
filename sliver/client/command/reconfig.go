package sliver_cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func RenameImplant(client *c.SliverClient, sessionID string, beaconID string, implantName string) (string, error) {
	session, _ := client.GetSessionByID(sessionID)
	beacon, _ := client.GetBeaconByID(beaconID)

	if session == nil && beacon == nil {
		return "", errors.New("need to provide either session or beacon")
	}

	_, err := client.RPC.Rename(context.Background(), &clientpb.RenameReq{
		SessionID: sessionID,
		BeaconID:  beaconID,
		Name:      implantName,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Renamed implant to %s", implantName), nil
}

// todo: add fn to change reconnect/beacon interval and beacon jitter
