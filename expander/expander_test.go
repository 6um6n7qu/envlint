package expander_test

import (
	"os"
	"testing"

	"github.com/user/envlint/expander"
)

func TestExpand_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	out := expander.Expand(env)
	if out["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected 5432, got %s", out["PORT"])
	}
}

func TestExpand_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "8080",
	}
	out := expander.Expand(env)
	want := "http://example.com:8080"
	if out["BASE_URL"] != want {
		t.Errorf("expected %q, got %q", want, out["BASE_URL"])
	}
}

func TestExpand_NoBraceStyle(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello $NAME",
		"NAME":     "World",
	}
	out := expander.Expand(env)
	if out["GREETING"] != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", out["GREETING"])
	}
}

func TestExpand_FallsBackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{
		"VALUE": "${OS_VAR}",
	}
	out := expander.Expand(env)
	if out["VALUE"] != "from-os" {
		t.Errorf("expected 'from-os', got %q", out["VALUE"])
	}
}

func TestExpand_UnresolvableBecomesEmpty(t *testing.T) {
	env := map[string]string{
		"VALUE": "prefix_${UNDEFINED_XYZ}_suffix",
	}
	out := expander.Expand(env)
	if out["VALUE"] != "prefix__suffix" {
		t.Errorf("expected 'prefix__suffix', got %q", out["VALUE"])
	}
}

func TestExpand_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "hello",
	}
	original := env["A"]
	expander.Expand(env)
	if env["A"] != original {
		t.Error("Expand must not mutate the input map")
	}
}
