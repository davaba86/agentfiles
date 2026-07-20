package templates

import (
	_ "embed"
	"fmt"
)

//go:embed README.md.tmpl
var ReadmeMD string

//go:embed Taskfile.yml.tmpl
var TaskfileYML string

//go:embed agentfiles.yml.tmpl
var AgentfilesYML string

// ClaudeMDBridge generates a CLAUDE.md bridge that references agentsPath.
func ClaudeMDBridge(agentsPath string) string {
	return fmt.Sprintf("@%s\n\n# Claude-specific notes\n\nShared project rules belong in AGENTS.md.\nOnly put Claude Code-specific behavior here.\n", agentsPath)
}

// ClaudeMD is the default bridge (AGENTS.md at repo root). Kept for backward compatibility.
var ClaudeMD = ClaudeMDBridge("AGENTS.md")
