package sliver_cmd

import (
	"context"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func Ps(client *c.SliverClient, session *clientpb.Session) (interface{}, error) {
	resp, err := client.RPC.Ps(context.Background(), &sliverpb.PsReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Terminate(client *c.SliverClient, session *clientpb.Session, pid int32) (interface{}, error) {
	resp, err := client.RPC.Terminate(context.Background(), &sliverpb.TerminateReq{
		Pid:     pid,
		Force:   true,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
