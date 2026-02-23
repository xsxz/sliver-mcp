package sliver_cmd

import (
	"context"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"

	//	"github.com/bishopfox/sliver/server/encoders"
	c "github.com/xsxz/sliver-mcp/sliver/client"
)

//func Upload(client *c.SliverClient, session *clientpb.Session, path string) (interface{}, error) {
//	resp, err := client.RPC.Upload(context.Background(), &sliverpb.UploadReq{
//		Path:    path,
//		Request: client.MakeSessionRequest(session),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return resp, nil
//}
//
//func Download(client *c.SliverClient, session *clientpb.Session, path string) (*sliverpb.Download, error) {
//	resp, err := client.RPC.Download(context.Background(), &sliverpb.DownloadReq{
//		Path:    path,
//		Request: client.MakeSessionRequest(session),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return resp, nil
//}

////
////
////

//func Cat(client *c.SliverClient, session *clientpb.Session, path string) (string, error) {
//	download, err := Download(client, session, path, true)
//	if err != nil {
//		return "", err
//	}
//
//	if download.Response != nil && download.Response.Err != "" {
//		return "", fmt.Errorf("download error: %s", download.Response.Err)
//	}
//
//	if download.Encoder == "gzip" {
//		download.Data, err = encoders.Gzip{}.Decode(download.Data)
//		if err != nil {
//			return "", fmt.Errorf("gzip decoding failed: %s", err)
//		}
//	}
//
//	return string(download.Data), nil
//}

func Ls(client *c.SliverClient, session *clientpb.Session, path string) (interface{}, error) {
	resp, err := client.RPC.Ls(context.Background(), &sliverpb.LsReq{
		Path:    path,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Cd(client *c.SliverClient, session *clientpb.Session, path string) (interface{}, error) {
	resp, err := client.RPC.Cd(context.Background(), &sliverpb.CdReq{
		Path:    path,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Pwd(client *c.SliverClient, session *clientpb.Session) (interface{}, error) {
	resp, err := client.RPC.Pwd(context.Background(), &sliverpb.PwdReq{
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Rm(client *c.SliverClient, session *clientpb.Session, path string) (interface{}, error) {
	resp, err := client.RPC.Rm(context.Background(), &sliverpb.RmReq{
		Path:    path,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Mv(client *c.SliverClient, session *clientpb.Session, src string, dst string) (interface{}, error) {
	resp, err := client.RPC.Mv(context.Background(), &sliverpb.MvReq{
		Src:     src,
		Dst:     dst,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Mkdir(client *c.SliverClient, session *clientpb.Session, path string) (interface{}, error) {
	resp, err := client.RPC.Mkdir(context.Background(), &sliverpb.MkdirReq{
		Path:    path,
		Request: client.MakeSessionRequest(session),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

////
////
////

func Chmod(client *c.SliverClient, session *clientpb.Session, filePath string, fileMode string, recursive bool) (*sliverpb.Chmod, error) {
	chmod, err := client.RPC.Chmod(context.Background(), &sliverpb.ChmodReq{
		Request:   client.MakeSessionRequest(session),
		Path:      filePath,
		FileMode:  fileMode,
		Recursive: recursive,
	})
	if err != nil {
		return nil, err
	}

	return chmod, nil
}

func Chown(client *c.SliverClient, session *clientpb.Session, filePath string, uid string, gid string, recursive bool) (*sliverpb.Chown, error) {
	chown, err := client.RPC.Chown(context.Background(), &sliverpb.ChownReq{
		Request:   client.MakeSessionRequest(session),
		Path:      filePath,
		Uid:       uid,
		Gid:       gid,
		Recursive: recursive,
	})
	if err != nil {
		return nil, err
	}

	return chown, nil
}

func Chtimes(client *c.SliverClient, session *clientpb.Session, filePath string, unixAtime int64, unixMtime int64) (*sliverpb.Chtimes, error) {
	chtimes, err := client.RPC.Chtimes(context.Background(), &sliverpb.ChtimesReq{
		Request: client.MakeSessionRequest(session),
		Path:    filePath,
		ATime:   unixAtime,
		MTime:   unixMtime,
	})
	if err != nil {
		return nil, err
	}

	return chtimes, nil
}
