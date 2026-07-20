package templates

import (
	"strings"
	"testing"
)

func TestClaudeMDBridgeContainsPath(t *testing.T) {
	result := ClaudeMDBridge(".agent-config/AGENTS.md")
	if !strings.Contains(result, "@.agent-config/AGENTS.md") {
		t.Fatalf("ClaudeMDBridge result missing expected path reference:\n%s", result)
	}
}

func TestClaudeMDBridgeDefaultMatchesClaudeMD(t *testing.T) {
	if ClaudeMDBridge("AGENTS.md") != ClaudeMD {
		t.Fatalf("ClaudeMDBridge(\"AGENTS.md\") must equal ClaudeMD for backward compatibility")
	}
}
