package masker_test

import (
	"flag"
	"testing"

	"github.com/yourorg/envlint/masker"
)

func TestFlagSet_DefaultsMaskEnabled(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts := masker.FlagSet(fs)

	if err := fs.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	if !opts.MaskSecrets {
		t.Error("expected MaskSecrets to default to true")
	}
}

func TestFlagSet_DisableMask(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	opts := masker.FlagSet(fs)

	if err := fs.Parse([]string{"-mask-secrets=false"}); err != nil {
		t.Fatal(err)
	}
	if opts.MaskSecrets {
		t.Error("expected MaskSecrets to be false")
	}
}

func TestOptions_Build_MaskEnabled(t *testing.T) {
	opts := &masker.Options{MaskSecrets: true}
	m := opts.Build()

	if !m.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be sensitive when masking enabled")
	}
}

func TestOptions_Build_MaskDisabled(t *testing.T) {
	opts := &masker.Options{MaskSecrets: false}
	m := opts.Build()

	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != "supersecret" {
		t.Errorf("expected unmasked value when masking disabled, got %q", got)
	}
}

func TestOptions_Build_ExtraPatterns(t *testing.T) {
	opts := &masker.Options{
		MaskSecrets:   true,
		ExtraPatterns: []string{"WEBHOOK_URL"},
	}
	m := opts.Build()

	if !m.IsSensitive("SLACK_WEBHOOK_URL") {
		t.Error("expected extra pattern WEBHOOK_URL to match SLACK_WEBHOOK_URL")
	}
	if !m.IsSensitive("DB_PASSWORD") {
		t.Error("expected default patterns to still apply alongside extra patterns")
	}
}
