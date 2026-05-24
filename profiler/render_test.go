package profiler_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlint/profiler"
)

func TestPrintTo_ContainsSections(t *testing.T) {
	r := profiler.Profile(makeSchema())
	var buf bytes.Buffer
	profiler.PrintTo(&buf, r)
	out := buf.String()

	for _, want := range []string{"Total", "Required", "Optional", "By group", "database", "app"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestPrintTo_EmptySchema(t *testing.T) {
	r := profiler.Report{ByGroup: map[string]int{}}
	var buf bytes.Buffer
	profiler.PrintTo(&buf, r)
	out := buf.String()
	if strings.Contains(out, "By group") {
		t.Error("expected no 'By group' section for empty schema")
	}
}

func TestPrintJSONTo_ValidJSON(t *testing.T) {
	r := profiler.Profile(makeSchema())
	var buf bytes.Buffer
	if err := profiler.PrintJSONTo(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded profiler.Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if decoded.Total != r.Total {
		t.Errorf("expected Total=%d, got %d", r.Total, decoded.Total)
	}
}

func TestPrintJSONTo_GroupsPreserved(t *testing.T) {
	r := profiler.Profile(makeSchema())
	var buf bytes.Buffer
	_ = profiler.PrintJSONTo(&buf, r)

	if !strings.Contains(buf.String(), "database") {
		t.Error("expected JSON to contain group 'database'")
	}
}
