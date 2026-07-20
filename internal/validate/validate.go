package validate

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/davaba86/agentfiles/internal/config"
	"github.com/davaba86/agentfiles/internal/files"
)

type Check struct {
	Message string
	OK      bool
	Warning bool
}

type Result struct {
	Checks []Check
}

func (r Result) OK() bool {
	for _, check := range r.Checks {
		if !check.OK && !check.Warning {
			return false
		}
	}
	return true
}

func Run(dir string, cfg config.Config, out io.Writer) (Result, error) {
	result := Result{}

	agentsPath := files.Join(dir, cfg.AgentsPath)
	claudePath := files.Join(dir, "CLAUDE.md")

	agentsExists, err := files.Exists(agentsPath)
	if err != nil {
		return result, fmt.Errorf("check %s: %w", cfg.AgentsPath, err)
	}
	result.Checks = append(result.Checks, foundCheck(cfg.AgentsPath, agentsExists))

	claudeExists, err := files.Exists(claudePath)
	if err != nil {
		return result, fmt.Errorf("check CLAUDE.md: %w", err)
	}
	result.Checks = append(result.Checks, foundCheck("CLAUDE.md", claudeExists))

	var agents, claude []byte
	if agentsExists {
		agents, err = files.Read(agentsPath)
		if err != nil {
			return result, fmt.Errorf("read %s: %w", cfg.AgentsPath, err)
		}
	}
	if claudeExists {
		claude, err = files.Read(claudePath)
		if err != nil {
			return result, fmt.Errorf("read CLAUDE.md: %w", err)
		}
	}

	if claudeExists {
		result.Checks = append(result.Checks, Check{
			Message: fmt.Sprintf("CLAUDE.md references %s", cfg.AgentsPath),
			OK:      strings.Contains(string(claude), cfg.AgentsPath),
		})
	}
	if agentsExists && claudeExists {
		result.Checks = append(result.Checks, Check{
			Message: fmt.Sprintf("CLAUDE.md is not a full duplicate of %s", cfg.AgentsPath),
			OK:      !bytes.Equal(bytes.TrimSpace(agents), bytes.TrimSpace(claude)),
		})
	}

	for _, name := range []string{"README.md", "Taskfile.yml", ".agentfiles.yml"} {
		exists, err := files.Exists(files.Join(dir, name))
		if err != nil {
			return result, fmt.Errorf("check %s: %w", name, err)
		}
		result.Checks = append(result.Checks, foundCheck(name, exists))
	}

	fmt.Fprintln(out, "Agentfiles check")
	fmt.Fprintln(out)
	for _, check := range result.Checks {
		status := "ok"
		if !check.OK {
			status = "fail"
		}
		fmt.Fprintf(out, "[%s] %s\n", status, check.Message)
	}
	fmt.Fprintln(out)
	if result.OK() {
		fmt.Fprintln(out, "Status: ok")
	} else {
		fmt.Fprintln(out, "Status: failed")
	}

	return result, nil
}

func foundCheck(name string, ok bool) Check {
	message := name + " found"
	if !ok {
		message = name + " missing"
	}
	return Check{Message: message, OK: ok}
}
