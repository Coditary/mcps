package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/coditary/mcps/pkg/config"
	"github.com/coditary/mcps/pkg/exec"
	"github.com/coditary/mcps/pkg/output"
	"github.com/coditary/mcps/pkg/schema"
	"github.com/coditary/mcps/pkg/template"
)

const confirmField = "confirm"

// RunConfig starts an MCP server from declarative config bytes.
func RunConfig(ctx context.Context, data []byte) error {
	cfg, err := config.Parse(data)
	if err != nil {
		return err
	}
	return Run(ctx, cfg)
}

// Run starts an MCP server from a parsed config.
func Run(ctx context.Context, cfg *config.Config) error {
	binary, err := exec.ResolveBinary(cfg.MCP.Binary, cfg.MCP.BinaryEnv)
	if err != nil {
		return err
	}
	timeout, err := exec.ParseTimeout(cfg.MCP.Timeout)
	if err != nil {
		return err
	}

	runner := &exec.Runner{
		Binary:  binary,
		WorkDir: cfg.MCP.WorkDir,
		Timeout: timeout,
	}

	title := cfg.MCP.Title
	if title == "" {
		title = cfg.MCP.Name
	}

	srv := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.MCP.Name,
		Title:   title,
		Version: cfg.MCP.Version,
	}, nil)

	for _, tool := range cfg.Tools {
		if err := registerTool(srv, runner, tool); err != nil {
			return fmt.Errorf("register tool %q: %w", tool.Name, err)
		}
	}

	return srv.Run(ctx, &mcp.StdioTransport{})
}

func registerTool(srv *mcp.Server, runner *exec.Runner, tool config.ToolSpec) error {
	inputSchema, err := schema.BuildInputSchemaJSON(tool)
	if err != nil {
		return err
	}

	srv.AddTool(&mcp.Tool{
		Name:        tool.Name,
		Description: tool.Description,
		InputSchema: inputSchema,
	}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleTool(ctx, runner, tool, req)
	})
	return nil
}

func handleTool(ctx context.Context, runner *exec.Runner, tool config.ToolSpec, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := map[string]any{}
	if len(req.Params.Arguments) > 0 {
		if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
			return nil, fmt.Errorf("decode arguments: %w", err)
		}
	}

	applyDefaults(tool, args)
	if err := validateRequired(tool, args); err != nil {
		return nil, err
	}

	confirm := false
	if tool.Mutating {
		if value, ok := args[confirmField]; ok {
			confirm, _ = value.(bool)
		}
		if !confirm {
			if tool.DryRunFlag == "" {
				return nil, fmt.Errorf("mutating tool requires confirm=true or a dry_run_flag")
			}
		}
		delete(args, confirmField)
	}

	stringArgs, err := template.Stringify(args)
	if err != nil {
		return nil, err
	}

	argv := expandSplitInputs(template.RenderArgs(tool.Exec.Args, stringArgs), tool.SplitInputs, stringArgs)
	for _, when := range tool.Exec.WhenArgs {
		if stringArgs[when.When] != "" {
			argv = append(argv, expandSplitInputs(template.RenderArgs(when.Args, stringArgs), tool.SplitInputs, stringArgs)...)
		}
	}
	if tool.Mutating && !confirm && tool.DryRunFlag != "" {
		argv = append(argv, tool.DryRunFlag)
	}

	result, err := runner.Run(ctx, argv, "")
	if err != nil {
		return nil, err
	}

	text, err := output.Format(tool.Output, stringArgs, result)
	if err != nil {
		return nil, err
	}

	isError := result.ExitCode != 0
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
		IsError: isError,
	}, nil
}

func applyDefaults(tool config.ToolSpec, args map[string]any) {
	for name, field := range tool.Inputs {
		if _, ok := args[name]; !ok && field.Default != nil {
			args[name] = field.Default
		}
	}
}

func expandSplitInputs(argv []string, splitInputs []string, stringArgs map[string]string) []string {
	if len(splitInputs) == 0 {
		return argv
	}

	splitSet := make(map[string]struct{}, len(splitInputs))
	for _, name := range splitInputs {
		splitSet[name] = struct{}{}
	}

	out := make([]string, 0, len(argv))
	for _, arg := range argv {
		expanded := false
		for name := range splitSet {
			value := strings.TrimSpace(stringArgs[name])
			if value == "" || arg != value {
				continue
			}
			out = append(out, strings.Fields(value)...)
			expanded = true
			break
		}
		if !expanded {
			out = append(out, arg)
		}
	}
	return out
}

func validateRequired(tool config.ToolSpec, args map[string]any) error {
	for name, field := range tool.Inputs {
		if !field.Required {
			continue
		}
		value, ok := args[name]
		if !ok || value == nil || value == "" {
			return fmt.Errorf("missing required argument %q", name)
		}
	}
	return nil
}

// RunFile loads config from disk and starts the server.
func RunFile(ctx context.Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return RunConfig(ctx, data)
}
