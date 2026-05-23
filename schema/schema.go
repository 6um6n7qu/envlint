package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// VarDefinition describes a single environment variable in the schema.
type VarDefinition struct {
	Name          string   `yaml:"name"`
	Required      bool     `yaml:"required"`
	Pattern       string   `yaml:"pattern"`
	AllowedValues []string `yaml:"allowed_values"`
	Default       string   `yaml:"default"`
	Description   string   `yaml:"description"`
	Group         string   `yaml:"group"`
}

// Schema is the top-level structure of a schema YAML file.
type Schema struct {
	Version string          `yaml:"version"`
	Vars    []VarDefinition `yaml:"vars"`
}

// Load reads and parses a schema YAML file from the given path.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read %q: %w", path, err)
	}
	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse %q: %w", path, err)
	}
	if len(s.Vars) == 0 {
		return nil, fmt.Errorf("schema: %q defines no variables", path)
	}
	return &s, nil
}
