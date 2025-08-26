package utils

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// https://github.com/ice-wzl/Sliver-Clients/commit/6a4f8b3035a6692c4bd9ac74ef7e7d89491edbaa#diff-bcd17627fea94208a39bb7e91fc98ee94c1fb4c721dc6677c379b3fcdadbec22R309
func FormatDateDelta(t time.Time, includeDate bool) string {
	nextTime := t.Format(time.UnixDate)

	var interval string

	if t.Before(time.Now()) {
		if includeDate {
			interval = fmt.Sprintf("%s (%s ago)", nextTime, time.Since(t).Round(time.Second))
		} else {
			interval = time.Since(t).Round(time.Second).String()
		}
	} else {
		if includeDate {
			interval = fmt.Sprintf("%s (in %s)", nextTime, time.Until(t).Round(time.Second))
		} else {
			interval = time.Until(t).Round(time.Second).String()
		}
	}
	return interval
}

// convert v1 protobuf proto.Message to JSON string
func ProtoMsgToJSON(msg proto.Message) (string, error) {
	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{}
	if err := marshaler.Marshal(&buf, msg); err != nil {
		return "", fmt.Errorf("failed to marshal protobuf message: %w", err)
	}
	return buf.String(), nil
}
