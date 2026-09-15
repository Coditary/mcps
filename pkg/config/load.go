package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and validates an MCP config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	return Parse(data)
}

// Parse unmarshals YAML config bytes.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate checks required fields.
func (c *Config) Validate() error {
	if c.MCP.Name == "" {
		return fmt.Errorf("mcp.name is required")
	}
	if c.MCP.Binary == "" {
		return fmt.Errorf("mcp.binary is required")
	}
	if c.MCP.Version == "" {
		c.MCP.Version = "v0.1.0"
	}
	for _, tool := range c.Tools {
		if tool.Name == "" {
			return fmt.Errorf("tool name is required")
		}
		if tool.Description == "" {
			return fmt.Errorf("tool %q: description is required", tool.Name)
		}
	}
	return nil
}
