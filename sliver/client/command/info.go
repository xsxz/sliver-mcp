package sliver_cmd

import (
	"time"

	c "github.com/xsxz/sliver-mcp/sliver/client"
	"github.com/xsxz/sliver-mcp/utils"
)

func GetSessionInfo(client *c.SliverClient, sessionID string) (map[string]interface{}, error) {
	session, err := client.GetSessionByID(sessionID)
	if err != nil {
		return nil, err //fmt.Errorf("execution failed: %w", err)
	}

	sessionDict := map[string]interface{}{
		"Session ID":         session.ID,
		"Implant Name":       session.Name,
		"Hostname":           session.Hostname,
		"UUID":               session.UUID,
		"Username":           session.Username,
		"UID":                session.UID,
		"GID":                session.GID,
		"Arch":               session.Arch,
		"OS":                 session.OS,
		"Locale":             session.Locale,
		"Transport":          session.Transport,
		"Remote Address":     session.RemoteAddress,
		"Proxy URL":          session.ProxyURL,
		"PID":                session.PID,
		"Filename":           session.Filename,
		"First Contact":      utils.FormatDateDelta(time.Unix(session.FirstContact, 0), true),
		"Last Checkin":       utils.FormatDateDelta(time.Unix(session.LastCheckin, 0), true),
		"Reconnect Interval": time.Duration(session.ReconnectInterval).String(),
		"Active C2":          session.ActiveC2,
		"Version":            session.Version,
	}

	return sessionDict, nil
}

func GetBeaconInfo(client *c.SliverClient, beaconID string) (map[string]interface{}, error) {
	beacon, err := client.GetBeaconByID(beaconID)
	if err != nil {
		return nil, err
	}

	beaconDict := map[string]interface{}{
		"Beacon ID":          beacon.ID,
		"Implant Name":       beacon.Name,
		"Hostname":           beacon.Hostname,
		"UUID":               beacon.UUID,
		"Username":           beacon.Username,
		"UID":                beacon.UID,
		"GID":                beacon.GID,
		"Arch":               beacon.Arch,
		"OS":                 beacon.OS,
		"Locale":             beacon.Locale,
		"Transport":          beacon.Transport,
		"Remote Address":     beacon.RemoteAddress,
		"Proxy URL":          beacon.ProxyURL,
		"PID":                beacon.PID,
		"Filename":           beacon.Filename,
		"Last Checkin":       utils.FormatDateDelta(time.Unix(beacon.LastCheckin, 0), true),
		"Next Checkin":       utils.FormatDateDelta(time.Unix(beacon.NextCheckin, 0), true),
		"First Contact":      utils.FormatDateDelta(time.Unix(beacon.FirstContact, 0), true),
		"Active C2":          beacon.ActiveC2,
		"Version":            beacon.Version,
		"Reconnect Interval": time.Duration(beacon.ReconnectInterval).String(),
		"Interval":           time.Duration(beacon.Interval).String(),
		"Jitter":             time.Duration(beacon.Jitter).String(),
		"Burned":             beacon.Burned,
		"Is Dead":            beacon.IsDead,
	}
	return beaconDict, nil
}
