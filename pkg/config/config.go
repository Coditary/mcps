package config

// Config describes an MCP server backed by a CLI binary.
type Config struct {
	MCP   MCPMeta    `yaml:"mcp"`
	Tools []ToolSpec `yaml:"tools"`
}

// MCPMeta holds server and binary resolution settings.
type MCPMeta struct {
	Name      string `yaml:"name"`
	Title     string `yaml:"title"`
	Version   string `yaml:"version"`
	Binary    string `yaml:"binary"`
	BinaryEnv string `yaml:"binary_env"`
	WorkDir   string `yaml:"work_dir"`
	Timeout   string `yaml:"timeout"`
}

// ToolSpec maps one MCP tool to a CLI invocation.
type ToolSpec struct {
	Name        string                `yaml:"name"`
	Description string                `yaml:"description"`
	Readonly    bool                  `yaml:"readonly"`
	Mutating    bool                  `yaml:"mutating"`
	DryRunFlag  string                `yaml:"dry_run_flag"`
	Exec        ExecSpec              `yaml:"exec"`
	SplitInputs []string              `yaml:"split_inputs"`
	Inputs      map[string]InputField `yaml:"inputs"`
	Output      OutputSpec            `yaml:"output"`
}

// ExecSpec describes argv passed to the wrapped binary.
type ExecSpec struct {
	Args     []string   `yaml:"args"`
	WhenArgs []WhenArgs `yaml:"when_args"`
}

// WhenArgs appends args when a named input is non-empty.
type WhenArgs struct {
	When string   `yaml:"when"`
	Args []string `yaml:"args"`
}

// InputField defines one tool parameter.
type InputField struct {
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Required    bool     `yaml:"required"`
	Default     any      `yaml:"default"`
	Enum        []string `yaml:"enum"`
}

// OutputSpec controls how CLI output is returned to the MCP client.
type OutputSpec struct {
	Kind    string `yaml:"kind"`
	Path    string `yaml:"path"`
	Pattern string `yaml:"pattern"`
}
