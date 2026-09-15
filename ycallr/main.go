package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/coditary/mcps/ycallr/apis"
)

func main() {
	store := &apis.Store{Dir: apis.DefaultDir()}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "ycallr",
		Title:   "ycallr MCP",
		Version: "v0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_apis",
		Description: "List configured ycallr API profiles.",
	}, listAPIs(store))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "describe_api",
		Description: "Describe commands and parameters for one ycallr API profile.",
	}, describeAPI(store))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "call_command",
		Description: "Execute a ycallr API command and return the JSON response body.",
	}, callCommand())

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "ycallr mcp failed: %v\n", err)
		os.Exit(1)
	}
}

type listAPIsInput struct{}

func listAPIs(store *apis.Store) func(context.Context, *mcp.CallToolRequest, listAPIsInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ listAPIsInput) (*mcp.CallToolResult, any, error) {
		items, err := store.List()
		if err != nil {
			return nil, nil, err
		}
		text, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return textResult(string(text)), nil, nil
	}
}

type describeAPIInput struct {
	API string `json:"api" jsonschema:"ycallr API profile name"`
}

func describeAPI(store *apis.Store) func(context.Context, *mcp.CallToolRequest, describeAPIInput) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, input describeAPIInput) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(input.API) == "" {
			return nil, nil, fmt.Errorf("api is required")
		}
		def, err := store.Load(input.API)
		if err != nil {
			return nil, nil, err
		}
		text, err := json.MarshalIndent(def, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return textResult(string(text)), nil, nil
	}
}

type callCommandInput struct {
	API     string            `json:"api" jsonschema:"ycallr API profile name"`
	Command string            `json:"command" jsonschema:"command name inside the API profile"`
	Params  map[string]string `json:"params,omitempty" jsonschema:"command parameters as key/value pairs"`
}

func callCommand() func(context.Context, *mcp.CallToolRequest, callCommandInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input callCommandInput) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(input.API) == "" || strings.TrimSpace(input.Command) == "" {
			return nil, nil, fmt.Errorf("api and command are required")
		}

		binary := "ycallr"
		if value := os.Getenv("YCALLR_BIN"); value != "" {
			binary = value
		}

		argv := []string{input.API, input.Command, "--json"}
		for key, value := range input.Params {
			argv = append(argv, fmt.Sprintf("--%s=%s", key, value))
		}

		cmd := exec.CommandContext(ctx, binary, argv...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: string(out)}},
				IsError: true,
			}, nil, nil
		}
		return textResult(string(out)), nil, nil
	}
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
