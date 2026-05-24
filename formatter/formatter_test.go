package formatter

import (
	"strings"
	"testing"
)

func makeEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestRender_Dotenv(t *testing.T) {
	out, err := Render(makeEntries(), FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected dotenv line, got: %s", out)
	}
}

func TestRender_Export(t *testing.T) {
	out, err := Render(makeEntries(), FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
	if !strings.Contains(out, "export DB_HOST=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestRender_JSON(t *testing.T) {
	out, err := Render(makeEntries(), FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"APP_ENV\"") {
		t.Errorf("expected JSON key APP_ENV, got: %s", out)
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestRender_UnknownFormat(t *testing.T) {
	_, err := Render(makeEntries(), Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "unknown format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRender_EmptyEntries(t *testing.T) {
	out, err := Render([]Entry{}, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}

func TestRender_Dotenv_ValueWithEquals(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Value: "postgres://user:pass@host/db?sslmode=disable"},
	}
	out, err := Render(entries, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DATABASE_URL=postgres://user:pass@host/db?sslmode=disable") {
		t.Errorf("expected full value preserved, got: %s", out)
	}
}
