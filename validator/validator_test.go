package validator_test

import (
	"testing"

	"github.com/envlint/schema"
	"github.com/envlint/validator"
)

func makeSchema(vars []schema.VarDef) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestValidate_AllPresent(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "PORT", Required: true},
		{Name: "HOST", Required: true},
	})
	env := map[string]string{"PORT": "8080", "HOST": "localhost"}
	report := validator.Validate(s, env)
	if report.HasErrors() {
		t.Fatalf("expected no errors, got: %v", report.Errors())
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "DATABASE_URL", Required: true},
	})
	env := map[string]string{}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Fatal("expected errors for missing required var")
	}
	if len(report.Errors()) != 1 || report.Errors()[0].Key != "DATABASE_URL" {
		t.Fatalf("unexpected errors: %v", report.Errors())
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "LOG_LEVEL", Required: false},
	})
	env := map[string]string{}
	report := validator.Validate(s, env)
	if report.HasErrors() {
		t.Fatalf("expected no errors for missing optional var, got: %v", report.Errors())
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "PORT", Required: true, Pattern: `^\d+$`},
	})
	env := map[string]string{"PORT": "8080"}
	report := validator.Validate(s, env)
	if report.HasErrors() {
		t.Fatalf("expected no errors, got: %v", report.Errors())
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "PORT", Required: true, Pattern: `^\d+$`},
	})
	env := map[string]string{"PORT": "not-a-number"}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Fatal("expected pattern mismatch error")
	}
}

func TestValidate_EmptyValueTreatedAsMissing(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "SECRET", Required: true},
	})
	env := map[string]string{"SECRET": "   "}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Fatal("expected error for blank required var")
	}
}
