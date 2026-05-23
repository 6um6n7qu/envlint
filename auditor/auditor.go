// Package auditor compares two env files and reports drift between them.
package auditor

import "fmt"

// DriftEntry describes a single variable that differs between two env sets.
type DriftEntry struct {
	Key      string
	BaseVal  string
	OtherVal string
	Status   string // "missing_in_other", "missing_in_base", "value_differs"
}

// Report holds the full drift comparison result.
type Report struct {
	Entries []DriftEntry
	HasDrift bool
}

// Compare takes two maps (e.g. staging vs production env vars) and returns
// a Report describing any drift between them.
func Compare(base, other map[string]string) Report {
	seen := make(map[string]bool)
	var entries []DriftEntry

	for k, baseVal := range base {
		seen[k] = true
		otherVal, ok := other[k]
		if !ok {
			entries = append(entries, DriftEntry{
				Key:     k,
				BaseVal: baseVal,
				Status:  "missing_in_other",
			})
		} else if baseVal != otherVal {
			entries = append(entries, DriftEntry{
				Key:      k,
				BaseVal:  baseVal,
				OtherVal: otherVal,
				Status:   "value_differs",
			})
		}
	}

	for k, otherVal := range other {
		if !seen[k] {
			entries = append(entries, DriftEntry{
				Key:      k,
				OtherVal: otherVal,
				Status:   "missing_in_base",
			})
		}
	}

	return Report{
		Entries:  entries,
		HasDrift: len(entries) > 0,
	}
}

// Summary returns a human-readable summary line for the report.
func Summary(r Report) string {
	if !r.HasDrift {
		return "No drift detected between env files."
	}
	return fmt.Sprintf("%d drift issue(s) detected.", len(r.Entries))
}
