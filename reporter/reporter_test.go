package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envlint/envlint/validator"
)

func makeResults() []validator.Result {
	return []validator.Result{
		{Key: "APP_ENV", Valid: true, Message: ""},
		{Key: "DB_URL", Valid: false, Message: "required variable is missing"},
		{Key: "PORT", Valid: false, Message: "value does not match pattern"},
	}
}

func TestPrint_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	Print(makeResults(), &buf)
	out := buf.String()

	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
	if !strings.Contains(out, "DB_URL") {
		t.Error("expected DB_URL in output")
	}
	if !strings.Contains(out, "required variable is missing") {
		t.Error("expected error message in output")
	}
}

func TestPrint_SummaryFailed(t *testing.T) {
	var buf bytes.Buffer
	Print(makeResults(), &buf)
	out := buf.String()

	if !strings.Contains(out, "Failed: 2") {
		t.Errorf("expected 'Failed: 2' in output, got:\n%s", out)
	}
}

func TestPrint_AllPassed(t *testing.T) {
	var buf bytes.Buffer
	results := []validator.Result{
		{Key: "APP_ENV", Valid: true},
		{Key: "PORT", Valid: true},
	}
	Print(results, &buf)
	out := buf.String()

	if !strings.Contains(out, "All 2 checks passed") {
		t.Errorf("expected all-passed summary, got:\n%s", out)
	}
}

func TestPrintJSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	PrintJSON(makeResults(), &buf)
	out := buf.String()

	if !strings.Contains(out, `"key": "APP_ENV"`) {
		t.Error("expected APP_ENV key in JSON output")
	}
	if !strings.Contains(out, `"valid": false`) {
		t.Error("expected false validity in JSON output")
	}
}
