package sliver_cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func StartHTTPListener(client *c.SliverClient, lhost string, lport uint16) (string, error) {
	httpListener, err := client.RPC.StartHTTPListener(context.Background(), &clientpb.HTTPListenerReq{
		Host:   lhost,
		Port:   uint32(lport),
		Secure: false,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully started HTTP listener on %s:%d. Job ID: %d", lhost, lport, httpListener.JobID), nil
}

func StartMTLSListener(client *c.SliverClient, lhost string, lport uint16) (string, error) {
	mtlsListener, err := client.RPC.StartMTLSListener(context.Background(), &clientpb.MTLSListenerReq{
		Host: lhost,
		Port: uint32(lport),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully started MTLS listener on %s:%d. Job ID: %d", lhost, lport, mtlsListener.JobID), nil
}

func StartWGListener(client *c.SliverClient, lport uint16, nport uint16, keyExchangePort uint16) (string, error) {
	wgListener, err := client.RPC.StartWGListener(context.Background(), &clientpb.WGListenerReq{
		Port:    uint32(lport),
		NPort:   uint32(nport),
		KeyPort: uint32(keyExchangePort),
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully started WG listener. Job ID: %d", wgListener.JobID), nil
}

func ListListeners(client *c.SliverClient) (map[uint32]*clientpb.Job, error) {
	jobs, err := client.RPC.GetJobs(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, err
	}

	listenerMap := make(map[uint32]*clientpb.Job)
	for _, job := range jobs.Active {
		listenerMap[job.ID] = job
	}

	return listenerMap, nil
}

func KillListener(client *c.SliverClient, jobID uint32) (string, error) {
	// need to add safeguard for grpc listener
	jobKill, err := client.RPC.KillJob(context.Background(), &clientpb.KillJobReq{
		ID: jobID,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully killed job #%d\n", jobKill.ID), nil
}

func KillAllListeners(client *c.SliverClient) ([]string, error) {
	listeners, err := ListListeners(client)
	if err != nil {
		return nil, err
	}

	var killed []string

	for _, job := range listeners {
		if strings.ToLower(job.Name) == "grpc" { // skip to avoid killing multiplayer mode listener
			continue
		}

		_, err := KillListener(client, job.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to kill listener %d: %w", job.ID, err)
		}

		killed = append(killed, fmt.Sprintf("ID: %d, Name: %s", job.ID, job.Name))
	}

	return killed, nil
}
