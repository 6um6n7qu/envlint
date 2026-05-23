package differ_test

import (
	"testing"

	"github.com/yourorg/envlint/differ"
)

func TestCompare_NoDiff(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "x", "B": "y"}

	d := differ.Compare(left, right)

	if !d.IsClean() {
		t.Errorf("expected clean diff, got: %+v", d)
	}
	if len(d.InBoth) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(d.InBoth))
	}
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := map[string]string{"A": "1", "EXTRA": "x"}
	right := map[string]string{"A": "1"}

	d := differ.Compare(left, right)

	if d.IsClean() {
		t.Error("expected dirty diff")
	}
	if len(d.OnlyInLeft) != 1 || d.OnlyInLeft[0] != "EXTRA" {
		t.Errorf("unexpected OnlyInLeft: %v", d.OnlyInLeft)
	}
	if len(d.OnlyInRight) != 0 {
		t.Errorf("expected empty OnlyInRight, got %v", d.OnlyInRight)
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "NEW_VAR": "hello"}

	d := differ.Compare(left, right)

	if len(d.OnlyInRight) != 1 || d.OnlyInRight[0] != "NEW_VAR" {
		t.Errorf("unexpected OnlyInRight: %v", d.OnlyInRight)
	}
}

func TestCompare_BothSidesMissing(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "2", "C": "3"}

	d := differ.Compare(left, right)

	if len(d.OnlyInLeft) != 1 || d.OnlyInLeft[0] != "A" {
		t.Errorf("unexpected OnlyInLeft: %v", d.OnlyInLeft)
	}
	if len(d.OnlyInRight) != 1 || d.OnlyInRight[0] != "C" {
		t.Errorf("unexpected OnlyInRight: %v", d.OnlyInRight)
	}
	if len(d.InBoth) != 1 || d.InBoth[0] != "B" {
		t.Errorf("unexpected InBoth: %v", d.InBoth)
	}
}

func TestSummary_Clean(t *testing.T) {
	d := differ.Diff{}
	if d.Summary() != "no differences found" {
		t.Errorf("unexpected summary: %s", d.Summary())
	}
}

func TestSummary_Dirty(t *testing.T) {
	d := differ.Diff{
		OnlyInLeft:  []string{"X"},
		OnlyInRight: []string{"Y", "Z"},
	}
	s := d.Summary()
	if s == "" || s == "no differences found" {
		t.Errorf("expected non-empty dirty summary, got: %s", s)
	}
}
