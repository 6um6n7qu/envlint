package filter_test

import (
	"testing"

	"github.com/user/envlint/filter"
	"github.com/user/envlint/schema"
)

func makeVars() []schema.VarDefinition {
	return []schema.VarDefinition{
		{Name: "DB_HOST", Required: true, Group: "database"},
		{Name: "DB_PORT", Required: true, Group: "database"},
		{Name: "APP_ENV", Required: true, Group: "app"},
		{Name: "LOG_LEVEL", Required: false, Group: "app"},
		{Name: "CACHE_TTL", Required: false, Group: "cache"},
	}
}

func TestApply_NoFilter(t *testing.T) {
	vars := makeVars()
	result := filter.Apply(vars, filter.Options{})
	if len(result) != len(vars) {
		t.Fatalf("expected %d, got %d", len(vars), len(result))
	}
}

func TestApply_ByPrefix(t *testing.T) {
	result := filter.Apply(makeVars(), filter.Options{Prefix: "DB_"})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	for _, v := range result {
		if v.Name[:3] != "DB_" {
			t.Errorf("unexpected var %s", v.Name)
		}
	}
}

func TestApply_ByGroup(t *testing.T) {
	result := filter.Apply(makeVars(), filter.Options{Groups: []string{"database"}})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestApply_ByRequired(t *testing.T) {
	req := true
	result := filter.Apply(makeVars(), filter.Options{Required: &req})
	if len(result) != 3 {
		t.Fatalf("expected 3 required vars, got %d", len(result))
	}
}

func TestApply_ByOptional(t *testing.T) {
	req := false
	result := filter.Apply(makeVars(), filter.Options{Required: &req})
	if len(result) != 2 {
		t.Fatalf("expected 2 optional vars, got %d", len(result))
	}
}

func TestApply_CombinedGroupAndRequired(t *testing.T) {
	req := false
	result := filter.Apply(makeVars(), filter.Options{
		Groups:   []string{"app"},
		Required: &req,
	})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Name != "LOG_LEVEL" {
		t.Errorf("unexpected var %s", result[0].Name)
	}
}
