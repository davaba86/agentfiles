package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultReturnsAgentsMD(t *testing.T) {
	cfg := Default()
	if cfg.AgentsPath != "AGENTS.md" {
		t.Fatalf("Default().AgentsPath = %q, want %q", cfg.AgentsPath, "AGENTS.md")
	}
}

func TestLoadMissingFileReturnsDefault(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.AgentsPath != "AGENTS.md" {
		t.Fatalf("AgentsPath = %q, want default %q", cfg.AgentsPath, "AGENTS.md")
	}
}

func TestLoadReadsCanonicalField(t *testing.T) {
	dir := t.TempDir()
	content := "canonical: .agent-config/AGENTS.md\n"
	if err := os.WriteFile(filepath.Join(dir, ".agentfiles.yml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.AgentsPath != ".agent-config/AGENTS.md" {
		t.Fatalf("AgentsPath = %q, want %q", cfg.AgentsPath, ".agent-config/AGENTS.md")
	}
}

func TestLoadDefaultCanonicalReturnsDefault(t *testing.T) {
	dir := t.TempDir()
	content := "canonical: AGENTS.md\n"
	if err := os.WriteFile(filepath.Join(dir, ".agentfiles.yml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.AgentsPath != "AGENTS.md" {
		t.Fatalf("AgentsPath = %q, want %q", cfg.AgentsPath, "AGENTS.md")
	}
}
