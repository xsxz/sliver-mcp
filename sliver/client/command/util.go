package sliver_cmd

import (
	"errors"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

func checkIfOsIsWin(client *c.SliverClient, session *clientpb.Session) error {
	targetOS, err := client.GetSessionByID(session.ID)
	if err != nil {
		return err
	}

	if targetOS.OS != "windows" {
		return errors.New("registry operations can only target Windows")
	}

	return nil
}
