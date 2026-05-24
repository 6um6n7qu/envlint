package linter

import (
	"testing"

	"github.com/your-org/envlint/schema"
)

func makeSchema(vars []schema.Var) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestLint_NoIssues(t *testing.T) {
	s := makeSchema([]schema.Var{
		{Name: "PORT", Required: true, Description: "HTTP port"},
		{Name: "DEBUG", Required: false, Default: "false"},
	})
	issues := Lint(s)
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_PatternAndAllowedValues(t *testing.T) {
	s := makeSchema([]schema.Var{
		{
			Name:          "ENV",
			Pattern:       "^(dev|prod)$",
			AllowedValues: []string{"dev", "prod"},
		},
	})
	issues := Lint(s)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %q", issues[0].Severity)
	}
}

func TestLint_RequiredNoDescription(t *testing.T) {
	s := makeSchema([]schema.Var{
		{Name: "SECRET_KEY", Required: true},
	})
	issues := Lint(s)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Var != "SECRET_KEY" {
		t.Errorf("expected issue for SECRET_KEY, got %q", issues[0].Var)
	}
}

func TestLint_RequiredWithDefault(t *testing.T) {
	s := makeSchema([]schema.Var{
		{Name: "DB_HOST", Required: true, Default: "localhost", Description: "database host"},
	})
	issues := Lint(s)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %q", issues[0].Severity)
	}
}

func TestLint_MultipleIssuesSameVar(t *testing.T) {
	s := makeSchema([]schema.Var{
		// required + no description (warning) + default (error)
		{Name: "BAD_VAR", Required: true, Default: "x"},
	})
	issues := Lint(s)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestLintAndValidate_ReturnsBoth(t *testing.T) {
	s := makeSchema([]schema.Var{
		{Name: "API_KEY", Required: true, Description: "API key"},
	})
	env := map[string]string{}
	issues, results := LintAndValidate(s, env)
	if len(issues) != 0 {
		t.Errorf("expected 0 lint issues, got %d", len(issues))
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 validation result, got %d", len(results))
	}
	if results[0].Valid {
		t.Error("expected validation failure for missing required var")
	}
}
