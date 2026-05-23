package auditor_test

import (
	"testing"

	"github.com/yourorg/envlint/auditor"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":   "staging",
		"DB_HOST":   "db.staging.internal",
		"LOG_LEVEL": "debug",
	}
}

func TestCompare_NoDrift(t *testing.T) {
	base := baseEnv()
	other := map[string]string{
		"APP_ENV":   "staging",
		"DB_HOST":   "db.staging.internal",
		"LOG_LEVEL": "debug",
	}
	r := auditor.Compare(base, other)
	if r.HasDrift {
		t.Errorf("expected no drift, got %d entries", len(r.Entries))
	}
}

func TestCompare_MissingInOther(t *testing.T) {
	base := baseEnv()
	other := map[string]string{
		"APP_ENV": "staging",
		"DB_HOST": "db.staging.internal",
		// LOG_LEVEL absent
	}
	r := auditor.Compare(base, other)
	if !r.HasDrift {
		t.Fatal("expected drift")
	}
	found := false
	for _, e := range r.Entries {
		if e.Key == "LOG_LEVEL" && e.Status == "missing_in_other" {
			found = true
		}
	}
	if !found {
		t.Error("expected missing_in_other entry for LOG_LEVEL")
	}
}

func TestCompare_MissingInBase(t *testing.T) {
	base := baseEnv()
	other := map[string]string{
		"APP_ENV":   "staging",
		"DB_HOST":   "db.staging.internal",
		"LOG_LEVEL": "debug",
		"NEW_RELIC":  "token123",
	}
	r := auditor.Compare(base, other)
	if !r.HasDrift {
		t.Fatal("expected drift")
	}
	for _, e := range r.Entries {
		if e.Key == "NEW_RELIC" && e.Status == "missing_in_base" {
			return
		}
	}
	t.Error("expected missing_in_base entry for NEW_RELIC")
}

func TestCompare_ValueDiffers(t *testing.T) {
	base := baseEnv()
	other := map[string]string{
		"APP_ENV":   "production",
		"DB_HOST":   "db.staging.internal",
		"LOG_LEVEL": "debug",
	}
	r := auditor.Compare(base, other)
	if !r.HasDrift {
		t.Fatal("expected drift")
	}
	for _, e := range r.Entries {
		if e.Key == "APP_ENV" && e.Status == "value_differs" {
			if e.BaseVal != "staging" || e.OtherVal != "production" {
				t.Errorf("unexpected values: %+v", e)
			}
			return
		}
	}
	t.Error("expected value_differs entry for APP_ENV")
}

func TestSummary_NoDrift(t *testing.T) {
	r := auditor.Report{HasDrift: false}
	if s := auditor.Summary(r); s != "No drift detected between env files." {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	r := auditor.Report{
		HasDrift: true,
		Entries:  []auditor.DriftEntry{{Key: "X", Status: "value_differs"}},
	}
	if s := auditor.Summary(r); s != "1 drift issue(s) detected." {
		t.Errorf("unexpected summary: %s", s)
	}
}
