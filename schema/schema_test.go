package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlint/schema"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp schema: %v", err)
	}
	return path
}

func TestLoad_ValidSchema(t *testing.T) {
	content := `
vars:
  PORT:
    description: HTTP port
    type: int
    required: true
  DEBUG:
    description: Enable debug mode
    type: bool
    required: false
    default: "false"
`
	path := writeTempFile(t, content)
	s, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}
	port, ok := s.Vars["PORT"]
	if !ok {
		t.Fatal("expected PORT var to be defined")
	}
	if port.Type != schema.TypeInt {
		t.Errorf("expected PORT type int, got %s", port.Type)
	}
	if !port.Required {
		t.Error("expected PORT to be required")
	}
}

func TestLoad_EmptyVars(t *testing.T) {
	path := writeTempFile(t, "vars:\n")
	s, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Vars == nil {
		t.Error("expected Vars map to be initialized, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := schema.Load("/nonexistent/path/schema.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTempFile(t, "vars: [invalid: yaml: content")
	_, err := schema.Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
