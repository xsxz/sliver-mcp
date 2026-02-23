package sliver_cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/bishopfox/sliver/client/command/generate"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
	"github.com/xsxz/sliver-mcp/utils"
)

func GetBeaconImplantConfig(client *c.SliverClient, beaconID string) (*clientpb.ImplantConfig, error) {
	beacon, err := client.GetBeaconByID(beaconID)
	if err != nil {
		return nil, err
	}

	c2s := []*clientpb.ImplantC2{}
	c2s = append(c2s, &clientpb.ImplantC2{
		URL:      beacon.ActiveC2,
		Priority: uint32(0),
	})

	config := &clientpb.ImplantConfig{
		ID:                  beacon.ID,
		Name:                beacon.GetName(),
		GOOS:                beacon.GetOS(),
		GOARCH:              beacon.GetArch(),
		Debug:               false,
		IsBeacon:            true,
		BeaconInterval:      beacon.GetInterval(),
		BeaconJitter:        beacon.GetJitter(),
		Evasion:             beacon.GetEvasion(),
		MaxConnectionErrors: uint32(1000),
		ReconnectInterval:   int64(60),
		Format:              clientpb.OutputFormat_SHELLCODE,
		IsSharedLib:         true,
		C2:                  c2s,
	}

	return config, nil
}

func KillBeacon(client *c.SliverClient, beaconID string, force bool) (string, error) {
	beacon, err := client.GetBeaconByID(beaconID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve beacon: %w", err)
	}

	_, err = client.RPC.Kill(context.Background(), &sliverpb.KillReq{
		Request: &commonpb.Request{
			BeaconID: beacon.ID,
			Timeout:  int64(utils.DefaultTimeout),
		},
		Force: force,
	})
	if err != nil {
		return "", fmt.Errorf("failed to kill beacon: %w", err)
	}

	return fmt.Sprintf("Beacon %s was killed", beaconID), nil
}

// StartBeaconInteractiveSession - Beacon only command to open an interactive session
func StartBeaconInteractiveSession(client *c.SliverClient, beaconID string) (interface{}, error) {
	beacon, err := client.GetBeaconByID(beaconID)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	var mtlsC2 []*clientpb.ImplantC2
	var wgC2 []*clientpb.ImplantC2
	var httpC2 []*clientpb.ImplantC2
	var dnsC2 []*clientpb.ImplantC2
	var namedPipeC2 []*clientpb.ImplantC2
	var tcpPivotC2 []*clientpb.ImplantC2

	c2s := []*clientpb.ImplantC2{}

	// parse the current beacon's ActiveC2
	c2url, err := url.Parse(beacon.ActiveC2)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	switch c2url.Scheme {
	case "mtls":
		mtlsC2, err = generate.ParseMTLSc2(beacon.ActiveC2)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		c2s = append(c2s, mtlsC2...)
	case "wg":
		wgC2, err = generate.ParseWGc2(beacon.ActiveC2)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		c2s = append(c2s, wgC2...)
	case "https":
		fallthrough
	case "http":
		httpC2, err = generate.ParseHTTPc2(beacon.ActiveC2)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		c2s = append(c2s, httpC2...)
	case "dns":
		dnsC2, err = generate.ParseDNSc2(beacon.ActiveC2)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		c2s = append(c2s, dnsC2...)
	case "namedpipe":
		namedPipeC2, err = generate.ParseNamedPipec2(beacon.ActiveC2)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		c2s = append(c2s, namedPipeC2...)
	case "tcppivot":
		tcpPivotC2, err = generate.ParseTCPPivotc2(beacon.ActiveC2)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		c2s = append(c2s, tcpPivotC2...)
	default:
		return nil, fmt.Errorf("unsupported C2 scheme: %s", c2url.Scheme)
	}

	openSession := &sliverpb.OpenSession{
		Request: &commonpb.Request{
			BeaconID: beacon.ID,
			Async:    true,
			Timeout:  int64(utils.DefaultTimeout),
		},
		C2S:   []string{},
		Delay: int64(0),
	}
	for _, c2 := range c2s {
		openSession.C2S = append(openSession.C2S, c2.URL)
	}

	openSession, err = client.RPC.OpenSession(context.Background(), openSession)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if openSession.Response != nil && openSession.Response.Err != "" {
		return nil, fmt.Errorf("%s", openSession.Response.Err)
	}

	return openSession.Response, nil

}

func RemoveBeacon(client *c.SliverClient, beaconID string) (string, error) {
	beacon, err := client.GetBeaconByID(beaconID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve beacon: %w", err)
	}

	_, err = client.RPC.RmBeacon(context.Background(), &clientpb.Beacon{ID: beacon.ID})
	if err != nil {
		return "", fmt.Errorf("failed to remove beacon: %w", err)
	}

	return fmt.Sprintf("Beacon %s was removed", beaconID), nil
}

func PruneBeacons(client *c.SliverClient, duration string) (string, error) {
	pruneDuration, err := time.ParseDuration(duration)
	if err != nil {
		return "", fmt.Errorf("unable to parse duration: %w", err)
	}

	beacons, err := client.ListBeacons()
	if err != nil {
		return "", fmt.Errorf("failed to get beacons: %w", err)
	}

	beaconsToPrune := []*clientpb.Beacon{}
	var prunedBeacons []string

	for _, beacon := range beacons.Beacons {
		nextCheckin := time.Unix(beacon.NextCheckin, 0)
		if time.Now().Before(nextCheckin) {
			continue
		}
		delta := time.Since(nextCheckin)
		if pruneDuration <= delta {
			beaconsToPrune = append(beaconsToPrune, beacon)
		}
	}

	if len(beaconsToPrune) == 0 {
		return "no beacons to prune", nil
	}

	for _, beacon := range beaconsToPrune {
		_, err := RemoveBeacon(client, beacon.ID)
		if err != nil {
			continue
		}
		prunedBeacons = append(prunedBeacons, beacon.ID)
	}

	if len(prunedBeacons) == 0 {
		return "no beacons were pruned", nil
	}

	result := fmt.Sprintf("Pruned beacons: %s", strings.Join(prunedBeacons, " "))
	return result, nil
}
