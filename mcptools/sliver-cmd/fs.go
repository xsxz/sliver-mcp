package mcptools_sliver_cmd

import (
	"context"
	"fmt"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/golang/protobuf/proto"
	"github.com/mark3labs/mcp-go/mcp"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
	cmd "github.com/xsxz/sliver-mcp/sliver/client/command"
	"github.com/xsxz/sliver-mcp/utils"
)

const (
	//Cat   = "session_run_cat"
	Ls    = "session_run_ls"
	Mv    = "session_run_mv"
	Cd    = "session_run_cd"
	Pwd   = "session_run_pwd"
	Rm    = "session_run_rm"
	Mkdir = "session_run_mkdir"

	Chmod   = "session_run_chmod"
	Chown   = "session_run_chown"
	Chtimes = "session_run_chtimes"
)

//var CatTool = mcp.NewTool(
//	Cat,
//	mcp.WithDescription("by session ID"),
//	mcp.WithString("session_id", mcp.Required()),
//	mcp.WithString("path", mcp.Required()),
//)
//
//func CatHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
//	path := request.GetString("path", "")
//	sessionID := request.GetString("session_id", "")
//	if sessionID == "" {
//		return nil, fmt.Errorf("session_id is required")
//	}
//
//	client, err := sliver_client.GetClient()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get sliver client: %w", err)
//	}
//
//	cmd := func(session *clientpb.Session) (interface{}, error) {
//		return cmd.Cat(client, session, path)
//	}
//
//	result, err := client.RunInSession(sessionID, cmd)
//	if err != nil {
//		return nil, fmt.Errorf("failed to run command: %w", err)
//	}
//
//	return &mcp.CallToolResult{
//		Content: []mcp.Content{
//			mcp.NewTextContent(result),
//		},
//	}, nil
//}

var LsTool = mcp.NewTool(
	Ls,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("path", mcp.Required()),
)

func LsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.GetString("path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Ls(client, session, path)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var CdTool = mcp.NewTool(
	Cd,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("path", mcp.Required()),
)

func CdHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.GetString("path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Cd(client, session, path)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var PwdTool = mcp.NewTool(
	Pwd,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
)

func PwdHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Pwd(client, session)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var MvTool = mcp.NewTool(
	Mv,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("source", mcp.Required()),
	mcp.WithString("destination", mcp.Required()),
)

func MvHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	src := request.GetString("source", "")
	dst := request.GetString("destination", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Mv(client, session, src, dst)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var RmTool = mcp.NewTool(
	Rm,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("path", mcp.Required()),
)

func RmHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.GetString("path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Rm(client, session, path)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var MkdirTool = mcp.NewTool(
	Mkdir,
	mcp.WithDescription("by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("path", mcp.Required()),
)

func MkdirHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := request.GetString("path", "")
	sessionID := request.GetString("session_id", "")
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Mkdir(client, session, path)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

////
////
////

var ChmodTool = mcp.NewTool(
	Chmod,
	mcp.WithDescription("Run chmod, by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("file_path", mcp.Required()),
	mcp.WithString("file_mode", mcp.Required()),
	mcp.WithBoolean("recursive"),
)

func ChmodHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	filePath := request.GetString("file_path", "")
	fileMode := request.GetString("file_mode", "")
	recursive := request.GetBool("recursive", false)

	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Chmod(client, session, filePath, fileMode, recursive)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run chmod command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chmod response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var ChownTool = mcp.NewTool(
	Chown,
	mcp.WithDescription("Run chown, by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("file_path", mcp.Required()),
	mcp.WithString("uid", mcp.Required()),
	mcp.WithString("gid", mcp.Required()),
	mcp.WithBoolean("recursive"),
)

func ChownHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	filePath := request.GetString("file_path", "")
	uid := request.GetString("uid", "")
	gid := request.GetString("gid", "")
	recursive := request.GetBool("recursive", false)

	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Chown(client, session, filePath, uid, gid, recursive)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run chown command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chown response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}

var ChtimesTool = mcp.NewTool(
	Chtimes,
	mcp.WithDescription("Run chtimes, by session ID"),
	mcp.WithString("session_id", mcp.Required()),
	mcp.WithString("file_path", mcp.Required()),
	mcp.WithNumber("unix_atime", mcp.Required()),
	mcp.WithNumber("unix_mtime", mcp.Required()),
)

func ChtimesHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	filePath := request.GetString("filePath", "")
	unixAtime := int64(request.GetInt("atime", 0))
	unixMtime := int64(request.GetInt("mtime", 0))

	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	client, err := sliver_client.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sliver client: %w", err)
	}

	cmd := func(session *clientpb.Session) (interface{}, error) {
		return cmd.Chtimes(client, session, filePath, unixAtime, unixMtime)
	}

	result, err := client.RunInSession(sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run chtimes command: %w", err)
	}

	jsonStrings, err := utils.ProtoMsgToJSON(result.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chtimes response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(jsonStrings),
		},
	}, nil
}
