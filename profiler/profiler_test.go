package profiler_test

import (
	"testing"

	"github.com/user/envlint/profiler"
	"github.com/user/envlint/schema"
)

func makeSchema() *schema.Schema {
	return &schema.Schema{
		Vars: []schema.Var{
			{Name: "DB_HOST", Required: true, Group: "database"},
			{Name: "DB_PASS", Required: true, Group: "database", Pattern: `^.{8,}$`},
			{Name: "LOG_LEVEL", Required: false, Default: "info", Allowed: []string{"debug", "info", "warn"}, Group: "app"},
			{Name: "PORT", Required: false, Default: "8080", Group: "app"},
			{Name: "SECRET_KEY", Required: true},
		},
	}
}

func TestProfile_Totals(t *testing.T) {
	r := profiler.Profile(makeSchema())
	if r.Total != 5 {
		t.Errorf("expected Total=5, got %d", r.Total)
	}
}

func TestProfile_RequiredOptional(t *testing.T) {
	r := profiler.Profile(makeSchema())
	if r.Required != 3 {
		t.Errorf("expected Required=3, got %d", r.Required)
	}
	if r.Optional != 2 {
		t.Errorf("expected Optional=2, got %d", r.Optional)
	}
}

func TestProfile_WithDefault(t *testing.T) {
	r := profiler.Profile(makeSchema())
	if r.WithDefault != 2 {
		t.Errorf("expected WithDefault=2, got %d", r.WithDefault)
	}
}

func TestProfile_WithPatternAndAllowed(t *testing.T) {
	r := profiler.Profile(makeSchema())
	if r.WithPattern != 1 {
		t.Errorf("expected WithPattern=1, got %d", r.WithPattern)
	}
	if r.WithAllowed != 1 {
		t.Errorf("expected WithAllowed=1, got %d", r.WithAllowed)
	}
}

func TestProfile_ByGroup(t *testing.T) {
	r := profiler.Profile(makeSchema())
	if r.ByGroup["database"] != 2 {
		t.Errorf("expected database=2, got %d", r.ByGroup["database"])
	}
	if r.ByGroup["app"] != 2 {
		t.Errorf("expected app=2, got %d", r.ByGroup["app"])
	}
	if r.ByGroup["(ungrouped)"] != 1 {
		t.Errorf("expected (ungrouped)=1, got %d", r.ByGroup["(ungrouped)"])
	}
}

func TestProfile_EmptySchema(t *testing.T) {
	r := profiler.Profile(&schema.Schema{})
	if r.Total != 0 {
		t.Errorf("expected Total=0, got %d", r.Total)
	}
}
