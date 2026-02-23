package sliver_cmd

import (
	"context"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func GetEnv(client *c.SliverClient, session *clientpb.Session, var_name string) (interface{}, error) {
	resp, err := client.RPC.GetEnv(context.Background(), &sliverpb.EnvReq{
		Name:    var_name,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SetEnv(client *c.SliverClient, session *clientpb.Session, var_name string, var_value string) (interface{}, error) {
	resp, err := client.RPC.SetEnv(context.Background(), &sliverpb.SetEnvReq{
		Variable: &commonpb.EnvVar{
			Key:   var_name,
			Value: var_value,
		},
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func UnsetEnv(client *c.SliverClient, session *clientpb.Session, var_name string) (interface{}, error) {
	resp, err := client.RPC.UnsetEnv(context.Background(), &sliverpb.UnsetEnvReq{
		Name:    var_name,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
