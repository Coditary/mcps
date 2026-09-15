package run

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/coditary/mcps/pkg/server"
)

// Main starts a declarative MCP server using an embedded config by default.
func Main(embedded []byte, defaultConfig string) {
	configPath := flag.String("config", "", "path to tools.yaml (overrides embedded config)")
	flag.Parse()

	ctx := context.Background()
	if *configPath != "" {
		if err := server.RunFile(ctx, *configPath); err != nil {
			fail(err)
		}
		return
	}

	data := embedded
	if len(data) == 0 && defaultConfig != "" {
		if err := server.RunFile(ctx, defaultConfig); err != nil {
			fail(err)
		}
		return
	}
	if len(data) == 0 {
		fmt.Fprintln(os.Stderr, "no config provided")
		os.Exit(2)
	}
	if err := server.RunConfig(ctx, data); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "mcp server failed: %v\n", err)
	os.Exit(1)
}
