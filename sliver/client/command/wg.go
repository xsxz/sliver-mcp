package sliver_cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func validateWG(session *clientpb.Session) error {
	if session.Transport != "wg" {
		return fmt.Errorf("this command is only supported for Wireguard implants")
	}
	return nil
}

func WGPortFwdList(client *c.SliverClient, session *clientpb.Session) (*sliverpb.WGTCPForwarders, error) {
	err := validateWG(session)
	if err != nil {
		return nil, err
	}

	fwdList, err := client.RPC.WGListForwarders(context.Background(), &sliverpb.WGTCPForwardersReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return fwdList, nil
}

func WGPortFwdAdd(client *c.SliverClient, session *clientpb.Session, localPort int, remoteAddr string) (*sliverpb.WGPortForward, error) {
	err := validateWG(session)
	if err != nil {
		return nil, err
	}

	if remoteAddr == "" {
		return nil, fmt.Errorf("must specify a remote target host:port")
	}

	//remoteHost, remotePort, err := net.SplitHostPort(remoteAddr)
	_, _, err = net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote target %s", err)

	}

	portfwdAdd, err := client.RPC.WGStartPortForward(context.Background(), &sliverpb.WGPortForwardStartReq{
		LocalPort:     int32(localPort),
		RemoteAddress: remoteAddr,
		Request:       client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return portfwdAdd, nil
}

func WGPortFwdRm(client *c.SliverClient, session *clientpb.Session, fwd_id int) (*sliverpb.WGPortForward, error) {
	err := validateWG(session)
	if err != nil {
		return nil, err
	}

	stopReq, err := client.RPC.WGStopPortForward(context.Background(), &sliverpb.WGPortForwardStopReq{
		ID:      int32(fwd_id),
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return stopReq, nil
}

///
///
///

func WGSocksList(client *c.SliverClient, session *clientpb.Session) (*sliverpb.WGSocksServers, error) {
	err := validateWG(session)
	if err != nil {
		return nil, err
	}

	socksList, err := client.RPC.WGListSocksServers(context.Background(), &sliverpb.WGSocksServersReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}
	if socksList.Response != nil && socksList.Response.Err != "" {
		return nil, fmt.Errorf("error: %s", socksList.Response.Err)
	}

	return socksList, nil

}
func WGSocksStart(client *c.SliverClient, session *clientpb.Session, bindPort int) (*sliverpb.WGSocks, error) {
	err := validateWG(session)
	if err != nil {
		return nil, err
	}

	socks, err := client.RPC.WGStartSocks(context.Background(), &sliverpb.WGSocksStartReq{
		Port:    int32(bindPort),
		Request: client.MakeSessionRequest(session),
	})

	if err != nil {
		return nil, err
	}
	if socks.Response != nil && socks.Response.Err != "" {
		return nil, fmt.Errorf("error: %s", socks.Response.Err)
	}

	return socks, nil

}
func WGSocksStop(client *c.SliverClient, session *clientpb.Session, socksID int) (*sliverpb.WGSocks, error) {
	err := validateWG(session)
	if err != nil {
		return nil, err
	}

	stopReq, err := client.RPC.WGStopSocks(context.Background(), &sliverpb.WGSocksStopReq{
		ID:      int32(socksID),
		Request: client.MakeSessionRequest(session),
	})

	if err != nil {
		return nil, err
	}
	if stopReq.Response != nil && stopReq.Response.Err != "" {
		return nil, fmt.Errorf("error: %s", stopReq.Response.Err)
	}

	return stopReq, nil

}
