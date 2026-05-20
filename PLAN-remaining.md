# Remaining Plan

## Deferred for the MVP

- Homebrew packaging and release automation.
- `migrate --force` support.
- Additional provider-specific behavior beyond the current shared convention.

## Still skeletal

- `internal/config/` exists but does not yet contain config parsing or validation logic.
- `.agentfiles.yml` is created by the tool, but the CLI does not read or act on it yet.

## Likely next steps

1. Add config loading in `internal/config` for `.agentfiles.yml`.
2. Teach `check` to honor config values such as `require_readme` and `require_taskfile`.
3. Expand `migrate` only after defining safe behavior for repos where both `AGENTS.md` and `CLAUDE.md` already exist.
4. Add release packaging once the CLI surface is stable.
