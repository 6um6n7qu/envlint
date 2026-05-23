package suggester_test

import (
	"strings"
	"testing"

	"github.com/user/envlint/suggester"
	"github.com/user/envlint/validator"
)

func makeResults(varName, message string, valid bool) validator.Result {
	return validator.Result{
		Variable: varName,
		Message:  message,
		Valid:    valid,
	}
}

func TestSuggest_SkipsValidResults(t *testing.T) {
	results := []validator.Result{
		makeResults("PORT", "ok", true),
		makeResults("HOST", "ok", true),
	}
	got := suggester.Suggest(results)
	if len(got) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(got))
	}
}

func TestSuggest_MissingRequired(t *testing.T) {
	results := []validator.Result{
		makeResults("DATABASE_URL", "missing required variable", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Hint, "required") {
		t.Errorf("expected hint to mention 'required', got: %s", got[0].Hint)
	}
	if got[0].Variable != "DATABASE_URL" {
		t.Errorf("expected variable DATABASE_URL, got %s", got[0].Variable)
	}
}

func TestSuggest_PatternMismatch(t *testing.T) {
	results := []validator.Result{
		makeResults("API_KEY", "value does not match pattern", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Hint, "format") {
		t.Errorf("expected hint to mention 'format', got: %s", got[0].Hint)
	}
}

func TestSuggest_AllowedValues(t *testing.T) {
	results := []validator.Result{
		makeResults("LOG_LEVEL", "value not in allowed list", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Hint, "allowed") {
		t.Errorf("expected hint to mention 'allowed', got: %s", got[0].Hint)
	}
}

func TestSuggest_FallbackHint(t *testing.T) {
	results := []validator.Result{
		makeResults("UNKNOWN_VAR", "some unexpected error", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Hint, "schema definition") {
		t.Errorf("expected fallback hint, got: %s", got[0].Hint)
	}
}
