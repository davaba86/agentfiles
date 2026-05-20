package migrate

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/davaba86/agentfiles/internal/files"
	"github.com/davaba86/agentfiles/templates"
)

type Options struct {
	DryRun bool
	Backup bool
}

func Run(dir string, opts Options, out io.Writer) error {
	claudePath := files.Join(dir, "CLAUDE.md")
	agentsPath := files.Join(dir, "AGENTS.md")
	claudeBackupPath := files.Join(dir, "CLAUDE.md.bak")
	agentsBackupPath := files.Join(dir, "AGENTS.md.bak")

	claudeExists, err := files.Exists(claudePath)
	if err != nil {
		return fmt.Errorf("check CLAUDE.md: %w", err)
	}
	agentsExists, err := files.Exists(agentsPath)
	if err != nil {
		return fmt.Errorf("check AGENTS.md: %w", err)
	}

	fmt.Fprintln(out, "Agentfiles migrate")
	fmt.Fprintln(out)

	switch {
	case !claudeExists && !agentsExists:
		fmt.Fprintln(out, "Neither CLAUDE.md nor AGENTS.md exists. Run `agentfiles init`.")
		return nil
	case claudeExists && agentsExists:
		return migrateExistingFiles(agentsPath, claudePath, agentsBackupPath, claudeBackupPath, opts, out)
	case !claudeExists && agentsExists:
		fmt.Fprintln(out, "AGENTS.md exists but CLAUDE.md is missing. Run `agentfiles init` to create the bridge.")
		return nil
	}

	fmt.Fprintln(out, "Will:")
	fmt.Fprintln(out, "  create AGENTS.md from CLAUDE.md")
	if opts.Backup {
		fmt.Fprintln(out, "  create CLAUDE.md.bak")
	}
	fmt.Fprintln(out, "  replace CLAUDE.md with the bridge template")

	if opts.DryRun {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Dry run: no files changed.")
		return nil
	}

	claude, err := files.Read(claudePath)
	if err != nil {
		return fmt.Errorf("read CLAUDE.md: %w", err)
	}
	if err := files.WriteNew(agentsPath, claude, 0o644); err != nil {
		return fmt.Errorf("create AGENTS.md: %w", err)
	}
	if opts.Backup {
		if err := files.WriteNew(claudeBackupPath, claude, 0o644); err != nil {
			return fmt.Errorf("create CLAUDE.md.bak: %w", err)
		}
	}
	if err := files.Write(claudePath, []byte(templates.ClaudeMD), 0o644); err != nil {
		return fmt.Errorf("replace CLAUDE.md: %w", err)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Migration complete.")
	return nil
}

func migrateExistingFiles(agentsPath, claudePath, agentsBackupPath, claudeBackupPath string, opts Options, out io.Writer) error {
	agents, err := files.Read(agentsPath)
	if err != nil {
		return fmt.Errorf("read AGENTS.md: %w", err)
	}
	claude, err := files.Read(claudePath)
	if err != nil {
		return fmt.Errorf("read CLAUDE.md: %w", err)
	}

	fmt.Fprintln(out, "Both AGENTS.md and CLAUDE.md exist.")
	fmt.Fprintln(out, "Will:")
	if shouldAppendClaudeRules(agents, claude) {
		fmt.Fprintln(out, "  append CLAUDE.md content to AGENTS.md")
	} else {
		fmt.Fprintln(out, "  keep AGENTS.md content unchanged")
	}
	if opts.Backup {
		fmt.Fprintln(out, "  create AGENTS.md.bak")
		fmt.Fprintln(out, "  create CLAUDE.md.bak")
	}
	fmt.Fprintln(out, "  replace CLAUDE.md with the bridge template")

	if opts.DryRun {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Dry run: no files changed.")
		return nil
	}

	if opts.Backup {
		if err := files.WriteNew(agentsBackupPath, agents, 0o644); err != nil {
			return fmt.Errorf("create AGENTS.md.bak: %w", err)
		}
		if err := files.WriteNew(claudeBackupPath, claude, 0o644); err != nil {
			return fmt.Errorf("create CLAUDE.md.bak: %w", err)
		}
	}
	if shouldAppendClaudeRules(agents, claude) {
		agents = appendClaudeRules(agents, claude)
		if err := files.Write(agentsPath, agents, 0o644); err != nil {
			return fmt.Errorf("update AGENTS.md: %w", err)
		}
	}
	if err := files.Write(claudePath, []byte(templates.ClaudeMD), 0o644); err != nil {
		return fmt.Errorf("replace CLAUDE.md: %w", err)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Migration complete.")
	return nil
}

func shouldAppendClaudeRules(agents, claude []byte) bool {
	claudeRules := bytes.TrimSpace(claude)
	if len(claudeRules) == 0 {
		return false
	}
	if bytes.Equal(claudeRules, bytes.TrimSpace([]byte(templates.ClaudeMD))) {
		return false
	}
	return !bytes.Contains(bytes.TrimSpace(agents), claudeRules)
}

func appendClaudeRules(agents, claude []byte) []byte {
	var out []byte
	out = append(out, bytes.TrimRight(agents, "\n")...)
	out = append(out, '\n', '\n')
	out = append(out, []byte("## Migrated from CLAUDE.md\n\n")...)
	out = append(out, []byte(strings.TrimSpace(string(claude)))...)
	out = append(out, '\n')
	return out
}
