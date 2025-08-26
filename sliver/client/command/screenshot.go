package sliver_cmd

///import (
///	"context"
///	"errors"
///	"fmt"
///	"time"
///
///	"github.com/bishopfox/sliver/protobuf/clientpb"
///	"github.com/bishopfox/sliver/protobuf/sliverpb"
///	c "github.com/xsxz/sliver-mcp/sliver/client"
///)

//func Screenshot(client *c.SliverClient, session *clientpb.Session) (string, error) {
//
//	if session.OS != "windows" && session.OS != "linux" {
//		return "", errors.New("target platform may not support screenshots")
//	}
//
//	screenshot, err := client.RPC.Screenshot(context.Background(), &sliverpb.ScreenshotReq{
//		Request: client.MakeSessionRequest(session),
//	})
//	if err != nil {
//		return "", err
//	}
//
//	if len(screenshot.Data) == 0 {
//		return "", errors.New("Cannot loot screenshot because it contained no data")
//	}
//
//	timeNow := time.Now().UTC()
//	screenshotFileName := fmt.Sprintf("screenshot_%s_%s.png", session.Hostname, timeNow.Format("20060102150405"))
//
//	//lootMessage := loot.CreateLootMessage(screenshotFileName, screenshotFileName, clientpb.LootType_LOOT_FILE, clientpb.FileType_BINARY, screenshot.GetData())
//	//loot.SendLootMessage(lootMessage, client)
//
//	//	LootScreenshot(screenshot, lootName, hostname, con)
//
//	//return "", nil
//}
