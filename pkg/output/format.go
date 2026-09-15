package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coditary/mcps/pkg/config"
	"github.com/coditary/mcps/pkg/exec"
	"github.com/coditary/mcps/pkg/template"
)

// Format builds the MCP tool response text from CLI output.
func Format(spec config.OutputSpec, vars map[string]string, result *exec.Result) (string, error) {
	switch spec.Kind {
	case "", "text":
		return formatText(result), nil
	case "json":
		return formatJSON(result)
	case "file":
		path := template.Render(spec.Path, vars)
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read output file %q: %w", path, err)
		}
		return string(data), nil
	case "files_glob":
		pattern := template.Render(spec.Pattern, vars)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return "", fmt.Errorf("glob %q: %w", pattern, err)
		}
		payload := map[string]any{
			"stdout":  result.Stdout,
			"stderr":  result.Stderr,
			"exit":    result.ExitCode,
			"files":   matches,
		}
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return formatText(result), nil
	}
}

func formatText(result *exec.Result) string {
	var b strings.Builder
	if result.Stdout != "" {
		b.WriteString(result.Stdout)
	}
	if result.Stderr != "" {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString("stderr:\n")
		b.WriteString(result.Stderr)
	}
	if result.ExitCode != 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("exit code: %d", result.ExitCode))
	}
	return b.String()
}

func formatJSON(result *exec.Result) (string, error) {
	text := strings.TrimSpace(result.Stdout)
	if text == "" {
		return formatText(result), nil
	}
	if !json.Valid([]byte(text)) {
		return formatText(result), nil
	}
	var parsed any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return formatText(result), nil
	}
	pretty, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return text, nil
	}
	return string(pretty), nil
}
