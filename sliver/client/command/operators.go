package sliver_cmd

import (
	"context"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func ListOperators(client *c.SliverClient) ([]*clientpb.Operator, error) {
	operatorsResp, err := client.RPC.GetOperators(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil, err
	}

	return operatorsResp.Operators, nil
}
