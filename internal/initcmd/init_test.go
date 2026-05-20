package initcmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/davaba86/agentfiles/internal/validate"
)

func TestRunCreatesMissingFiles(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	result, err := Run(dir, &out)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(result.Created) != 5 {
		t.Fatalf("created %d files, want 5", len(result.Created))
	}
	for _, name := range []string{"AGENTS.md", "CLAUDE.md", "README.md", "Taskfile.yml", ".agentfiles.yml"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected %s to exist: %v", name, err)
		}
	}
}

func TestRunCreatesRepoThatPassesValidation(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	if _, err := Run(dir, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	result, err := validate.Run(dir, &out)
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	if !result.OK() {
		t.Fatalf("expected initialized repo to pass validation")
	}
}

func TestRunDoesNotOverwriteExistingFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	const original = "keep me"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	if _, err := Run(dir, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != original {
		t.Fatalf("AGENTS.md was overwritten: got %q, want %q", data, original)
	}
}
