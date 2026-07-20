package validate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davaba86/agentfiles/internal/config"
)

func TestRunPassesForValidRepo(t *testing.T) {
	dir := validRepo(t)
	var out bytes.Buffer

	result, err := Run(dir, config.Default(), &out)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !result.OK() {
		t.Fatalf("expected valid repo to pass")
	}
}

func TestRunOutputsASCIIStatusMarkers(t *testing.T) {
	dir := validRepo(t)
	var out bytes.Buffer

	if _, err := Run(dir, config.Default(), &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "[ok] AGENTS.md found") {
		t.Fatalf("expected ASCII status marker in output, got:\n%s", output)
	}
	if strings.Contains(output, "✅") || strings.Contains(output, "❌") {
		t.Fatalf("output contains emoji status marker:\n%s", output)
	}
}

func TestRunFailsWhenAgentsMissing(t *testing.T) {
	dir := validRepo(t)
	if err := os.Remove(filepath.Join(dir, "AGENTS.md")); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	result, err := Run(dir, config.Default(), &out)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.OK() {
		t.Fatalf("expected missing AGENTS.md to fail")
	}
}

func TestRunFailsWhenClaudeDoesNotReferenceAgents(t *testing.T) {
	dir := validRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("no bridge here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	result, err := Run(dir, config.Default(), &out)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.OK() {
		t.Fatalf("expected missing AGENTS.md reference to fail")
	}
}

func TestRunPassesForSubmodulePath(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, ".agent-config")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	filemap := map[string]string{
		".agent-config/AGENTS.md": "# shared rules\n",
		"CLAUDE.md":               "@.agent-config/AGENTS.md\n",
		"README.md":               "# readme\n",
		"Taskfile.yml":            "version: '3'\n",
		".agentfiles.yml":         "canonical: .agent-config/AGENTS.md\n",
	}
	for name, content := range filemap {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cfg := config.Config{AgentsPath: ".agent-config/AGENTS.md"}
	var out bytes.Buffer
	result, err := Run(dir, cfg, &out)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !result.OK() {
		t.Fatalf("expected submodule repo to pass:\n%s", out.String())
	}
}

func validRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"AGENTS.md":       "# shared rules\n",
		"CLAUDE.md":       "@AGENTS.md\n",
		"README.md":       "# readme\n",
		"Taskfile.yml":    "version: '3'\n",
		".agentfiles.yml": "canonical: AGENTS.md\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}
