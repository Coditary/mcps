package apis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Definition is the subset of a ycallr API profile needed by the MCP server.
type Definition struct {
	Name        string             `yaml:"name"`
	Version     string             `yaml:"version"`
	Description string             `yaml:"description"`
	Commands    map[string]Command `yaml:"commands"`
}

// Command describes one API command from a profile.
type Command struct {
	Description string                 `yaml:"description"`
	Endpoint    string                 `yaml:"endpoint"`
	Method      string                 `yaml:"method"`
	Params      map[string]Param       `yaml:"params"`
}

// Param describes one command parameter.
type Param struct {
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
}

// Summary is returned by list_apis.
type Summary struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	CommandCount int   `json:"command_count"`
}

// Store loads API profiles from disk.
type Store struct {
	Dir string
}

// DefaultDir resolves the API profile directory.
func DefaultDir() string {
	if value := os.Getenv("YCALLR_APIS_DIR"); value != "" {
		return value
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "ycallr", "apis")
}

// Load reads one API profile by name.
func (s *Store) Load(name string) (*Definition, error) {
	path := filepath.Join(s.Dir, name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read api profile %q: %w", name, err)
	}
	var def Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parse api profile %q: %w", name, err)
	}
	if def.Name == "" {
		def.Name = name
	}
	return &def, nil
}

// List returns summaries for all API profiles in the store directory.
func (s *Store) List() ([]Summary, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, fmt.Errorf("read apis dir %q: %w", s.Dir, err)
	}

	out := make([]Summary, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		def, err := s.Load(name)
		if err != nil {
			continue
		}
		out = append(out, Summary{
			Name:         def.Name,
			Version:      def.Version,
			Description:  def.Description,
			CommandCount: len(def.Commands),
		})
	}
	return out, nil
}
