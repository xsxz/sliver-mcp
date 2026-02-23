package sliver_cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

// reg-create
// reg-delete
// reg-list
// reg-read
// reg-write

// used in create, deleted, read, write
func checkHive(hive string) error {
	if map[string]bool{"HKCU": true, "HKLM": true, "HKCC": true, "HKPD": true, "HKU": true, "HKCR": true}[hive] {
		return nil
	}
	return errors.New("invalid hive")
}

func getType(t string) (uint32, error) {
	var res uint32
	switch t {
	case "binary":
		res = sliverpb.RegistryTypeBinary
	case "dword":
		res = sliverpb.RegistryTypeDWORD
	case "qword":
		res = sliverpb.RegistryTypeQWORD
	case "string":
		res = sliverpb.RegistryTypeString
	default:
		return res, fmt.Errorf("invalid type %s", t)
	}
	return res, nil
}

// used in create, delete, read, write
func parseRegistryPath(regPath string) (finalPath, key string, err error) {
	if strings.Contains(regPath, "/") {
		regPath = strings.ReplaceAll(regPath, "/", "\\")
	}
	pathBaseIdx := strings.LastIndex(regPath, `\`)
	if pathBaseIdx < 0 {
		err = fmt.Errorf("invalid path: %s", regPath)
		return "", "", err
	}
	if len(regPath) < pathBaseIdx+1 {
		err = fmt.Errorf("invalid path: %s", regPath)
		return "", "", err
	}
	finalPath = regPath[:pathBaseIdx]
	key = regPath[pathBaseIdx+1:]
	return finalPath, key, nil
}

////
////
////

// List sub registry keys
func RegListSubKeys(client *c.SliverClient, session *clientpb.Session, hive string, regPath string) ([]string, error) {
	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
		return nil, err
	}

	subkeys, err := client.RPC.RegistryListSubKeys(context.Background(), &sliverpb.RegistrySubKeyListReq{
		Hive:     hive,
		Hostname: session.Hostname,
		Path:     regPath,
		Request:  client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return subkeys.Subkeys, nil
}

// List registry values
func RegListValues(client *c.SliverClient, session *clientpb.Session, hive string, regPath string) ([]string, error) {
	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
		return nil, err
	}

	values, err := client.RPC.RegistryListValues(context.Background(), &sliverpb.RegistryListValuesReq{
		Hive:     hive,
		Hostname: session.Hostname,
		Path:     regPath,
		Request:  client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return values.ValueNames, nil
}

func RegReadValue(client *c.SliverClient, session *clientpb.Session, hive string, regPath string) (string, error) {
	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
		return "", err
	}

	finalPath, key, err := parseRegistryPath(regPath)
	if err != nil {
		return "", err
	}

	regRead, err := client.RPC.RegistryRead(context.Background(), &sliverpb.RegistryReadReq{
		Hive:     hive,
		Path:     finalPath,
		Key:      key,
		Hostname: session.Hostname,
		Request:  client.MakeSessionRequest(session),
	})
	if err != nil {
		return "couldn't read from registry", err
	}

	return regRead.Value, nil
}

////
////
////

// List registry VALUES and then get KEYS for each of them
func RegReadKeysValues(client *c.SliverClient, session *clientpb.Session, hive string, regPath string) (string, error) {
	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
		return "", err
	}

	var results []string

	values, err := RegListValues(client, session, hive, regPath)
	if err != nil {
		return "", err
	}

	for _, value := range values {
		fullPath := regPath + "\\" + value

		keyValue, err := RegReadValue(client, session, hive, fullPath)
		if err != nil {
			continue
		}

		results = append(results, value+": "+keyValue)
	}

	return strings.Join(results, "; "), nil
}

func RegCreateKey(client *c.SliverClient, session *clientpb.Session, hive string, regPath string) (string, error) {
	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
	}

	finalPath, key, err := parseRegistryPath(regPath)
	if err != nil {
		return "", err
	}

	_, err = client.RPC.RegistryCreateKey(context.Background(), &sliverpb.RegistryCreateKeyReq{
		Hive:     hive,
		Path:     finalPath,
		Key:      key,
		Hostname: session.Hostname,
		Request:  client.MakeSessionRequest(session),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Key created at %s\\%s", regPath, key), nil
}

func RegWriteValue(client *c.SliverClient, session *clientpb.Session, hive string, regPath string, flagType string, value string) (string, error) {
	var (
		dwordValue  uint32
		qwordValue  uint64
		stringValue string
		binaryValue []byte
	)

	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
	}

	finalPath, key, err := parseRegistryPath(regPath)
	if err != nil {
		return "", err
	}

	valType, err := getType(flagType)

	switch valType {
	case sliverpb.RegistryTypeBinary:
		var (
			v   []byte
			err error
		)
		v, err = hex.DecodeString(value)
		if err != nil {
			return "", err
		}
		binaryValue = v
	case sliverpb.RegistryTypeDWORD:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return "", err
		}
		dwordValue = uint32(v)
	case sliverpb.RegistryTypeQWORD:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return "", err
		}
		qwordValue = v
	case sliverpb.RegistryTypeString:
		stringValue = value
	default:
		return "invalid type", err
	}

	_, err = client.RPC.RegistryWrite(context.Background(), &sliverpb.RegistryWriteReq{
		Request:     client.MakeSessionRequest(session),
		Hostname:    session.Hostname,
		Hive:        hive,
		Path:        finalPath,
		Type:        valType,
		Key:         key,
		StringValue: stringValue,
		DWordValue:  dwordValue,
		QWordValue:  qwordValue,
		ByteValue:   binaryValue,
	})
	if err != nil {
		return "", err
	}

	return "Value written to registry", nil
}

func RegDeleteKey(client *c.SliverClient, session *clientpb.Session, hive string, regPath string) (string, error) {
	checkIfOsIsWin(client, session)

	err := checkHive(hive)
	if err != nil {
	}

	finalPath, key, err := parseRegistryPath(regPath)
	if err != nil {
		return "", err
	}

	_, err = client.RPC.RegistryDeleteKey(context.Background(), &sliverpb.RegistryDeleteKeyReq{
		Hive:     hive,
		Path:     finalPath,
		Key:      key,
		Hostname: session.Hostname,
		Request:  client.MakeSessionRequest(session),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Key removed at %s\\%s", regPath, key), nil
}
