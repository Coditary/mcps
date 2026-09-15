package exec

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Result captures CLI execution output.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Runner executes wrapped CLI binaries.
type Runner struct {
	Binary  string
	WorkDir string
	Timeout time.Duration
}

// Run executes argv with optional stdin.
func (r *Runner) Run(ctx context.Context, argv []string, stdin string) (*Result, error) {
	if r.Binary == "" {
		return nil, fmt.Errorf("binary is not configured")
	}

	cmdCtx := ctx
	var cancel context.CancelFunc
	if r.Timeout > 0 {
		cmdCtx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(cmdCtx, r.Binary, argv...)
	if r.WorkDir != "" {
		cmd.Dir = r.WorkDir
	}
	if stdin != "" {
		cmd.Stdin = bytes.NewBufferString(stdin)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("run %s: %w", r.Binary, err)
		}
	}

	return &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

// ResolveBinary returns the binary path from config and optional env override.
func ResolveBinary(binary, envKey string) (string, error) {
	if envKey != "" {
		if value := os.Getenv(envKey); value != "" {
			return value, nil
		}
	}
	if binary == "" {
		return "", fmt.Errorf("binary is not configured")
	}
	return binary, nil
}

// ParseTimeout parses duration strings like "30s" or "2m".
func ParseTimeout(raw string) (time.Duration, error) {
	if raw == "" {
		return 5 * time.Minute, nil
	}
	return time.ParseDuration(raw)
}
