package migrate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/davaba86/agentfiles/templates"
)

func TestRunDryRunDoesNotWriteFiles(t *testing.T) {
	dir := oldClaudeRepo(t)
	var out bytes.Buffer

	if err := Run(dir, Options{DryRun: true, Backup: true}, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("AGENTS.md should not exist after dry run")
	}
	if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md.bak")); !os.IsNotExist(err) {
		t.Fatalf("CLAUDE.md.bak should not exist after dry run")
	}
}

func TestRunMigratesClaudeFirstRepo(t *testing.T) {
	dir := oldClaudeRepo(t)
	var out bytes.Buffer

	if err := Run(dir, Options{Backup: true}, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(agents) != "old claude rules\n" {
		t.Fatalf("AGENTS.md = %q, want copied Claude content", agents)
	}

	backup, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md.bak"))
	if err != nil {
		t.Fatal(err)
	}
	if string(backup) != "old claude rules\n" {
		t.Fatalf("CLAUDE.md.bak = %q, want original Claude content", backup)
	}

	claude, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(claude) != templates.ClaudeMD {
		t.Fatalf("CLAUDE.md = %q, want bridge template", claude)
	}
}

func TestRunMigratesRepoWithExistingAgents(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# Agent Instructions\n\nExisting rules.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Claude Instructions\n\nClaude rules.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	if err := Run(dir, Options{Backup: true}, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	wantAgents := "# Agent Instructions\n\nExisting rules.\n\n## Migrated from CLAUDE.md\n\n# Claude Instructions\n\nClaude rules.\n"
	if string(agents) != wantAgents {
		t.Fatalf("AGENTS.md = %q, want merged content", agents)
	}

	agentsBackup, err := os.ReadFile(filepath.Join(dir, "AGENTS.md.bak"))
	if err != nil {
		t.Fatal(err)
	}
	if string(agentsBackup) != "# Agent Instructions\n\nExisting rules.\n" {
		t.Fatalf("AGENTS.md.bak = %q, want original AGENTS.md content", agentsBackup)
	}

	claudeBackup, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md.bak"))
	if err != nil {
		t.Fatal(err)
	}
	if string(claudeBackup) != "# Claude Instructions\n\nClaude rules.\n" {
		t.Fatalf("CLAUDE.md.bak = %q, want original CLAUDE.md content", claudeBackup)
	}

	claude, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(claude) != templates.ClaudeMD {
		t.Fatalf("CLAUDE.md = %q, want bridge template", claude)
	}
}

func TestRunDoesNotAppendClaudeBridgeToExistingAgents(t *testing.T) {
	dir := t.TempDir()
	const agentsContent = "# Agent Instructions\n"
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte(agentsContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(templates.ClaudeMD), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	if err := Run(dir, Options{Backup: true}, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(agents) != agentsContent {
		t.Fatalf("AGENTS.md = %q, want unchanged bridge-only migration", agents)
	}

	claude, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(claude) != templates.ClaudeMD {
		t.Fatalf("CLAUDE.md = %q, want bridge template", claude)
	}
}

func TestRunDoesNotAppendClaudeRulesAlreadyInAgents(t *testing.T) {
	dir := t.TempDir()
	const agentsContent = "# Claude Instructions\n\nClaude rules.\n"
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte(agentsContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(agentsContent), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	if err := Run(dir, Options{Backup: true}, &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(agents) != agentsContent {
		t.Fatalf("AGENTS.md = %q, want unchanged existing rules", agents)
	}
}

func oldClaudeRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("old claude rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}
