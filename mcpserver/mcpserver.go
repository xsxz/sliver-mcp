package mcpserver

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/server"
	sc "github.com/xsxz/sliver-mcp/mcptools/sliver-client"
	scmd "github.com/xsxz/sliver-mcp/mcptools/sliver-cmd"

	"github.com/xsxz/sliver-mcp/utils"
)

var (
	version = utils.Version
)

func addTools(s *server.MCPServer) {
	//// client commands
	// listeners
	s.AddTool(sc.ListListenersTool, sc.ListListenersHandle)
	s.AddTool(sc.KillAllListenersTool, sc.KillAllListenersHandle)
	s.AddTool(sc.StartHTTPListenerTool, sc.StartHTTPListenerHandle)
	s.AddTool(sc.StartMTLSListenerTool, sc.StartMTLSListenerHandle)
	s.AddTool(sc.StartWGListenerTool, sc.StartWGListenerHandle)
	s.AddTool(sc.KillListenerTool, sc.KillListenerHandle)

	// sessions
	s.AddTool(sc.ListSessionsTool, sc.ListSessionsHandle)
	s.AddTool(sc.GetSessionInfoTool, sc.GetSessionInfoHandle)
	s.AddTool(sc.KillSessionTool, sc.KillSessionHandle)
	s.AddTool(sc.KillAllSessionsTool, sc.KillAllSessionsHandle)
	s.AddTool(sc.CloseInteractiveSessionTool, sc.CloseInteractiveSessionHandle)
	s.AddTool(sc.PruneSessionsTool, sc.PruneSessionsToolHandle)

	// beacons
	s.AddTool(sc.ListBeaconsTool, sc.ListBeaconsHandle)
	s.AddTool(sc.GetBeaconInfoTool, sc.GetBeaconInfoHandle)
	s.AddTool(sc.KillBeaconTool, sc.KillBeaconHandle)
	s.AddTool(sc.KillAllBeaconsTool, sc.KillAllBeaconsHandle)
	s.AddTool(sc.StartBeaconInteractiveSessionTool, sc.StartBeaconInteractiveSessionHandle)
	s.AddTool(sc.PruneBeaconsTool, sc.PruneBeaconsHandle)

	// operators
	s.AddTool(sc.ListOperatorsTool, sc.ListOperatorsHandle)

	// reconfig
	s.AddTool(sc.RenameImplantTool, sc.RenameImplantHandle)

	// generate
	s.AddTool(sc.ListImplantBuildsTool, sc.ListImplantBuildsHandle)
	s.AddTool(sc.RmImplantBuildTool, sc.RmImplantBuildHandle)

	s.AddTool(sc.ListImplantProfilesTool, sc.ListImplantProfilesHandle)
	s.AddTool(sc.RmImplantProfileTool, sc.RmImplantProfileHandle)
	s.AddTool(sc.GenerateWinShellcodeProfileTool, sc.GenerateWinShellcodeProfileHandle) // for demo purposes only

	//// session/beacon commands
	// pivots
	s.AddTool(scmd.ListPivotsTool, scmd.ListPivotsHandle)
	s.AddTool(scmd.StopAllPivotsTool, scmd.StopAllPivotsHandle)
	s.AddTool(scmd.StartNamedPipePivotTool, scmd.StartNamedPipePivotHandle)
	s.AddTool(scmd.StartTcpPivotTool, scmd.StartTcpPivotHandle)
	s.AddTool(scmd.StopPivotTool, scmd.StopPivotHandle)

	// wg

	s.AddTool(scmd.ListWGSocksTool, scmd.ListWGSocksHandle)
	s.AddTool(scmd.ListWGPortFwdTool, scmd.ListWGPortFwdHandle)
	s.AddTool(scmd.StartWGSocksTool, scmd.StartWGSocksHandle)
	s.AddTool(scmd.StartWGPortFwdTool, scmd.StartWGPortFwdHandle)
	s.AddTool(scmd.StopWGSocksTool, scmd.StopWGSocksHandle)
	s.AddTool(scmd.StopWGPortFwdTool, scmd.StopWGPortFwdHandle)

	// net
	s.AddTool(scmd.IfconfigTool, scmd.IfconfigHandle)
	s.AddTool(scmd.NetstatTool, scmd.NetstatHandle)

	// env
	s.AddTool(scmd.SetEnvTool, scmd.SetEnvHandle)
	s.AddTool(scmd.GetEnvTool, scmd.GetEnvHandle)
	s.AddTool(scmd.UnsetEnvTool, scmd.UnsetEnvHandle)

	// fs
	//s.AddTool(scmd.CatTool, scmd.CatHandle)
	s.AddTool(scmd.LsTool, scmd.LsHandle)
	s.AddTool(scmd.CdTool, scmd.CdHandle)
	s.AddTool(scmd.PwdTool, scmd.PwdHandle)
	s.AddTool(scmd.MvTool, scmd.MvHandle)
	s.AddTool(scmd.RmTool, scmd.RmHandle)
	s.AddTool(scmd.MkdirTool, scmd.MkdirHandle)
	s.AddTool(scmd.ChmodTool, scmd.ChmodHandle)
	s.AddTool(scmd.ChownTool, scmd.ChownHandle)
	s.AddTool(scmd.ChtimesTool, scmd.ChtimesHandle)

	// processes
	s.AddTool(scmd.PsTool, scmd.PsHandle)
	s.AddTool(scmd.TerminateTool, scmd.TerminateHandle)

	//exec
	s.AddTool(scmd.ExecuteTool, scmd.ExecuteHandle)
	s.AddTool(scmd.MsfTool, scmd.MsfHandle)
	s.AddTool(scmd.MsfInjectTool, scmd.MsfInjectHandle)
	s.AddTool(scmd.MigrateTool, scmd.MigrateHandle)

	// win-specific commands

	// registry
	s.AddTool(scmd.RegReadKeysValuesTool, scmd.RegReadKeysValuesHandle)
	s.AddTool(scmd.RegReadValueTool, scmd.RegReadValueHandle)
	s.AddTool(scmd.RegCreateKeyTool, scmd.RegCreateKeyHandle)
	s.AddTool(scmd.RegWriteValueTool, scmd.RegWriteValueHandle)
	s.AddTool(scmd.RegDeleteKeyTool, scmd.RegDeleteKeyHandle)

	//privs
	s.AddTool(scmd.GetPrivsTool, scmd.GetPrivsHandle)
	s.AddTool(scmd.GetSystemTool, scmd.GetSystemHandle)
	s.AddTool(scmd.ImpersonateTool, scmd.ImpersonateHandle)
	s.AddTool(scmd.MakeTokenTool, scmd.MakeTokenHandle)
	s.AddTool(scmd.Rev2SelfTool, scmd.Rev2SelfHandle)
	s.AddTool(scmd.RunAsTool, scmd.RunAsHandle)

	s.AddTool(scmd.BackdoorTool, scmd.BackdoorHandle)

}

func StartServer(listenIP string, listenPort int) {
	if listenIP == "" {
		listenIP = "localhost"
	}

	if listenPort == 0 {
		listenPort = 11337
	}
	address := fmt.Sprintf("%s:%d", listenIP, listenPort)
	slog.Info("Starting Sliver MCP server on " + address)

	s := server.NewMCPServer("Sliver MCP server", version,
		server.WithToolCapabilities(true),
		server.WithRecovery(), // automatically recover from panics
	)

	addTools(s)

	httpServer := server.NewStreamableHTTPServer(s,
		server.WithHeartbeatInterval(30*time.Second),
		server.WithStateLess(true),
	)

	if err := httpServer.Start(address); err != nil {
		log.Fatal(err)
	}

}
