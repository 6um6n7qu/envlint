package masker_test

import (
	"testing"

	"github.com/yourorg/envlint/masker"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := masker.New()

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"GITHUB_TOKEN",
		"AWS_SECRET",
		"PRIVATE_KEY_PATH",
		"USER_AUTH_TOKEN",
	}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	m := masker.New()

	safe := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
	}
	for _, key := range safe {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMask_SensitiveValue(t *testing.T) {
	m := masker.New()
	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != "***" {
		t.Errorf("expected masked value, got %q", got)
	}
}

func TestMask_SafeValue(t *testing.T) {
	m := masker.New()
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestMaskAll(t *testing.T) {
	m := masker.New()
	input := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"PORT":        "8080",
	}
	result := m.MaskAll(input)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked")
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be masked")
	}
	if result["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] != "***" {
		t.Errorf("API_KEY should be masked")
	}
}

func TestNewWithPatterns(t *testing.T) {
	m := masker.NewWithPatterns([]string{"CUSTOM_SECRET"})
	if !m.IsSensitive("MY_CUSTOM_SECRET_VAL") {
		t.Error("expected custom pattern to match")
	}
	if m.IsSensitive("DB_PASSWORD") {
		t.Error("default patterns should not apply when using custom patterns")
	}
}
