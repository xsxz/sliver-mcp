package sliver_cmd

import (
	"context"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func Ifconfig(client *c.SliverClient, session *clientpb.Session) (interface{}, error) {
	resp, err := client.RPC.Ifconfig(context.Background(), &sliverpb.IfconfigReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Netstat(client *c.SliverClient, session *clientpb.Session) (interface{}, error) {
	resp, err := client.RPC.Netstat(context.Background(), &sliverpb.NetstatReq{
		TCP:       true,
		UDP:       true,
		IP4:       true,
		IP6:       true,
		Listening: true,
		Request:   client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
