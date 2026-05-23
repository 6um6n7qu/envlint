package filter_test

import (
	"flag"
	"testing"

	"github.com/user/envlint/filter"
)

func TestFlagSet_Prefix(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	buildOpts := filter.FlagSet(fs)
	_ = fs.Parse([]string{"--prefix", "APP_"})
	opts := buildOpts()
	if opts.Prefix != "APP_" {
		t.Errorf("expected prefix APP_, got %q", opts.Prefix)
	}
}

func TestFlagSet_Groups(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	buildOpts := filter.FlagSet(fs)
	_ = fs.Parse([]string{"--group", "database,app"})
	opts := buildOpts()
	if len(opts.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(opts.Groups))
	}
}

func TestFlagSet_RequiredOnly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	buildOpts := filter.FlagSet(fs)
	_ = fs.Parse([]string{"--required"})
	opts := buildOpts()
	if opts.Required == nil || *opts.Required != true {
		t.Error("expected Required=true")
	}
}

func TestFlagSet_OptionalOnly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	buildOpts := filter.FlagSet(fs)
	_ = fs.Parse([]string{"--optional"})
	opts := buildOpts()
	if opts.Required == nil || *opts.Required != false {
		t.Error("expected Required=false")
	}
}

func TestFlagSet_NoFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	buildOpts := filter.FlagSet(fs)
	_ = fs.Parse([]string{})
	opts := buildOpts()
	if opts.Prefix != "" || len(opts.Groups) != 0 || opts.Required != nil {
		t.Error("expected zero-value options with no flags")
	}
}
