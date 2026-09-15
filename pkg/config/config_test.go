package config

import "testing"

func TestParseMinimalConfig(t *testing.T) {
	data := []byte(`
mcp:
  name: demo
  binary: demo-cli
tools:
  - name: ping
    description: ping the cli
    exec:
      args: ["--version"]
    inputs: {}
`)
	cfg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if cfg.MCP.Name != "demo" {
		t.Fatalf("expected demo, got %q", cfg.MCP.Name)
	}
	if len(cfg.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(cfg.Tools))
	}
}
