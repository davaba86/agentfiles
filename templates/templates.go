package templates

import _ "embed"

//go:embed AGENTS.md.tmpl
var AgentsMD string

//go:embed CLAUDE.md.tmpl
var ClaudeMD string

//go:embed README.md.tmpl
var ReadmeMD string

//go:embed Taskfile.yml.tmpl
var TaskfileYML string

//go:embed agentfiles.yml.tmpl
var AgentfilesYML string
