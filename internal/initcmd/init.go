package initcmd

import (
	"fmt"
	"io"

	"github.com/davaba86/agentfiles/internal/files"
	"github.com/davaba86/agentfiles/templates"
)

type Result struct {
	Created []string
	Skipped []string
}

func Run(dir string, out io.Writer) (Result, error) {
	items := []struct {
		name    string
		content string
	}{
		{name: "AGENTS.md", content: templates.AgentsMD},
		{name: "CLAUDE.md", content: templates.ClaudeMD},
		{name: "README.md", content: templates.ReadmeMD},
		{name: "Taskfile.yml", content: templates.TaskfileYML},
		{name: ".agentfiles.yml", content: templates.AgentfilesYML},
	}

	result := Result{}
	for _, item := range items {
		path := files.Join(dir, item.name)
		exists, err := files.Exists(path)
		if err != nil {
			return result, fmt.Errorf("check %s: %w", item.name, err)
		}
		if exists {
			result.Skipped = append(result.Skipped, item.name+" already exists")
			continue
		}
		if err := files.WriteNew(path, []byte(item.content), 0o644); err != nil {
			return result, fmt.Errorf("create %s: %w", item.name, err)
		}
		result.Created = append(result.Created, item.name)
	}

	fmt.Fprintln(out, "Agentfiles init")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Created:")
	for _, name := range result.Created {
		fmt.Fprintf(out, "  %s\n", name)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Skipped:")
	for _, name := range result.Skipped {
		fmt.Fprintf(out, "  %s\n", name)
	}

	return result, nil
}
