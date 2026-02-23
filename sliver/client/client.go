package sliver_client

// todo:

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/bishopfox/sliver/client/assets"
	"github.com/xsxz/sliver-mcp/utils"
	"google.golang.org/grpc"

	"github.com/bishopfox/sliver/client/transport"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
)

//// sliver client dependency injection

var (
	client *SliverClient
)

func SetClient(c *SliverClient) {
	client = c
}

func GetClient() (*SliverClient, error) {
	if client == nil {
		return nil, errors.New("sliver client is not initialized yet")
	}
	return client, nil
}

////
////
////

type SessionCommand func(session *clientpb.Session) (interface{}, error)
type BeaconCommand func(session *clientpb.Beacon) (interface{}, error)

type SliverClient struct {
	RPC         rpcpb.SliverRPCClient
	Conn        *grpc.ClientConn
	EventStream rpcpb.SliverRPC_EventsClient
}

func NewClient(configPath string) (*SliverClient, error) {
	config, err := assets.ReadConfig(configPath)
	if err != nil {
		return nil, errors.New("failed to read Sliver config")
	}

	rpc, conn, err := transport.MTLSConnect(config)
	if err != nil {
		return nil, errors.New("failed to connect to Sliver server")
	}
	slog.Info("Connected to Sliver server")

	eventStream, err := rpc.Events(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, errors.New("failed to open event stream")
	}

	client := &SliverClient{
		RPC:         rpc,
		Conn:        conn,
		EventStream: eventStream,
	}

	return client, nil
}

func (c *SliverClient) ListSessions() (*clientpb.Sessions, error) {
	sessions, err := c.RPC.GetSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, errors.New("failed to get sessions")
	}

	return sessions, nil
}

func (c *SliverClient) ListBeacons() (*clientpb.Beacons, error) {
	beacons, err := c.RPC.GetBeacons(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, errors.New("failed to get beacons")
	}

	return beacons, nil
}

func (c *SliverClient) GetSessionByID(sessionID string) (*clientpb.Session, error) {
	sessions, err := c.ListSessions()
	if err != nil {
		return nil, errors.New("session not found")
	}

	for _, session := range sessions.Sessions {
		if session.ID == sessionID {
			return session, nil
		}
	}

	return nil, errors.New("session not found")
}

func (c *SliverClient) GetBeaconByID(beaconID string) (*clientpb.Beacon, error) {
	beacons, err := c.ListBeacons()
	if err != nil {
		return nil, errors.New("failed to list beacons")
	}

	for _, beacon := range beacons.Beacons {
		if beacon.ID == beaconID {
			return beacon, nil
		}
	}

	return nil, errors.New("beacon not found")
}

////
////
////

func (c *SliverClient) MakeSessionRequest(session *clientpb.Session) *commonpb.Request {
	if session == nil {
		return nil
	}

	return &commonpb.Request{
		SessionID: session.ID,
		Timeout:   int64(utils.DefaultTimeout), // 60
	}
}

func (c *SliverClient) RunInSession(sessionID string, cmd SessionCommand) (interface{}, error) {
	session, err := c.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	return cmd(session)
}

func (c *SliverClient) RunInBeacon(beaconID string, cmd BeaconCommand) (interface{}, error) {
	beacon, err := c.GetBeaconByID(beaconID)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	return cmd(beacon)
}
