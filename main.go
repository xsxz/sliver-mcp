package main

import (
	"flag"
	"os"

	"log/slog"

	"github.com/xsxz/sliver-mcp/mcpserver"
	sliver_client "github.com/xsxz/sliver-mcp/sliver/client"
)

func main() {

	configPath := flag.String("operator-config", "", "Path to Sliver client config file (required)")
	lhost := flag.String("lhost", "", "Address to listen MCP HTTP server, default is 127.0.0.1")
	lport := flag.Int("lport", 0, "Port to listen MCP HTTP server, default is 11337")
	flag.Parse()

	if *configPath == "" {
		slog.Error("[!] Error: -config flag is required")
		flag.Usage()
		os.Exit(1)
	}

	sc, err := sliver_client.NewClient(*configPath)
	if err != nil {
		slog.Error("Failed to load sliver client", slog.Any("err", err))
	}
	defer sc.Conn.Close()

	sliver_client.SetClient(sc)

	mcpserver.StartServer(*lhost, *lport)

}
