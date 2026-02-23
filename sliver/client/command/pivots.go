package sliver_cmd

import (
	"context"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func ListPivots(client *c.SliverClient, session *clientpb.Session) (*sliverpb.PivotListeners, error) {
	pivotListeners, err := client.RPC.PivotSessionListeners(context.Background(), &sliverpb.PivotListenersReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return pivotListeners, nil
}

func StartTcpPivot(client *c.SliverClient, session *clientpb.Session, bindAddress string, lport uint16) (string, error) {
	listener, err := client.RPC.PivotStartListener(context.Background(), &sliverpb.PivotStartListenerReq{
		Request:     client.MakeSessionRequest(session),
		Type:        sliverpb.PivotType_TCP,
		BindAddress: fmt.Sprintf("%s:%d", bindAddress, lport),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Started TCP pivot listener %s with id %d\n", listener.BindAddress, listener.ID), nil
}

func StartNamedPipePivot(client *c.SliverClient, session *clientpb.Session, bindAddress string, allowAll bool) (string, error) {
	options := []bool{allowAll}

	listener, err := client.RPC.PivotStartListener(context.Background(), &sliverpb.PivotStartListenerReq{
		Request:     client.MakeSessionRequest(session),
		Type:        sliverpb.PivotType_NamedPipe,
		BindAddress: bindAddress,
		Options:     options,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Started named pipe pivot listener %s with id %d\n", listener.BindAddress, listener.ID), nil
}

func StopPivot(client *c.SliverClient, session *clientpb.Session, ID uint32) (string, error) {
	pivotListeners, err := ListPivots(client, session)
	if err != nil {
		return "", err
	}

	if len(pivotListeners.Listeners) == 0 {
		return "", fmt.Errorf("no pivot listeners running on this session")
	}

	_, err = client.RPC.PivotStopListener(context.Background(), &sliverpb.PivotStopListenerReq{
		Request: client.MakeSessionRequest(session),
		ID:      ID,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Stopped pivot listener with ID %d.\n", ID), nil
}

func StopAllPivots(client *c.SliverClient, session *clientpb.Session) (string, error) {
	pivotListeners, err := ListPivots(client, session)
	if err != nil {
		return "", err
	}

	if len(pivotListeners.Listeners) == 0 {
		return "", fmt.Errorf("no pivot listeners running on this session")
	}

	for _, listener := range pivotListeners.Listeners {
		_, err := client.RPC.PivotStopListener(context.Background(), &sliverpb.PivotStopListenerReq{
			Request: client.MakeSessionRequest(session),
			ID:      listener.ID,
		})
		if err != nil {
			return "", fmt.Errorf("failed to stop pivot listener with ID %d: %w", listener.ID, err)
		}
	}

	return "Stopped all pivot listeners.\n", nil
}
