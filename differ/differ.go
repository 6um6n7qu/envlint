// Package differ compares two sets of env variables and reports
// keys that are present in one but absent in the other.
package differ

import (
	"fmt"
	"sort"
)

// Diff represents the result of comparing two env maps.
type Diff struct {
	OnlyInLeft  []string
	OnlyInRight []string
	InBoth      []string
}

// Compare returns a Diff between left and right env maps.
// Keys are compared by name only; values are ignored.
func Compare(left, right map[string]string) Diff {
	d := Diff{}

	for k := range left {
		if _, ok := right[k]; ok {
			d.InBoth = append(d.InBoth, k)
		} else {
			d.OnlyInLeft = append(d.OnlyInLeft, k)
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			d.OnlyInRight = append(d.OnlyInRight, k)
		}
	}

	sort.Strings(d.OnlyInLeft)
	sort.Strings(d.OnlyInRight)
	sort.Strings(d.InBoth)

	return d
}

// IsClean returns true when both sides contain exactly the same keys.
func (d Diff) IsClean() bool {
	return len(d.OnlyInLeft) == 0 && len(d.OnlyInRight) == 0
}

// Summary returns a one-line human-readable description of the diff.
func (d Diff) Summary() string {
	if d.IsClean() {
		return "no differences found"
	}
	result := ""
	if len(d.OnlyInLeft) > 0 {
		result += fmt.Sprintf("%d key(s) only in left", len(d.OnlyInLeft))
	}
	if len(d.OnlyInRight) > 0 {
		if result != "" {
			result += "; "
		}
		result += fmt.Sprintf("%d key(s) only in right", len(d.OnlyInRight))
	}
	return result
}
