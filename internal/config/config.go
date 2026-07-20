package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AgentsPath string
}

func Default() Config {
	return Config{AgentsPath: "AGENTS.md"}
}

func Load(dir string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(filepath.Join(dir, ".agentfiles.yml"))
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if after, ok := strings.CutPrefix(line, "canonical:"); ok {
			if path := strings.TrimSpace(after); path != "" {
				cfg.AgentsPath = path
			}
		}
	}
	return cfg, nil
}
