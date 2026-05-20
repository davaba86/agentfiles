# agentfiles

`agentfiles` is a small CLI for standardizing AI coding-agent instruction files across repositories.

## Convention

- `AGENTS.md` is the canonical shared instruction file.
- `CLAUDE.md` is a minimal Claude Code bridge that references `AGENTS.md`.
- `README.md` is for human-facing project documentation.
- `Taskfile.yml` defines executable project workflows.
- `.agentfiles.yml` stores project configuration for this tool.

## Commands

```sh
go run ./cmd/agentfiles init
go run ./cmd/agentfiles check
go run ./cmd/agentfiles migrate [--dry-run]
go run ./cmd/agentfiles version
```

## Local development

Build the CLI:

```sh
task build
```

Then run the built binary from another repository by absolute or relative path:

```sh
cd /path/to/other/repo

# Safe migration preview.
../agentfiles/bin/agentfiles migrate --dry-run

# Apply the migration.
../agentfiles/bin/agentfiles migrate

# Validate the result.
../agentfiles/bin/agentfiles check
```

Use `migrate --dry-run` first when testing against an existing repository. Non-dry-run `migrate` writes files in the target repository.

Migration behavior:

- If only `CLAUDE.md` exists, `migrate` creates `AGENTS.md` from it, backs up `CLAUDE.md`, and replaces `CLAUDE.md` with the bridge template.
- If both `AGENTS.md` and `CLAUDE.md` exist, `migrate` appends substantive Claude rules into `AGENTS.md`, backs up both files, and replaces `CLAUDE.md` with the bridge template.
- If `CLAUDE.md` is already just the bridge, `migrate` leaves `AGENTS.md` content unchanged.

To create starter files in a new repository instead of migrating existing instructions:

```sh
../agentfiles/bin/agentfiles init
```

## Homebrew

The intended public install path is a Homebrew tap:

```sh
brew install davaba86/tap/agentfiles
```

The tap should live in a separate public repository named `homebrew-tap`. See [docs/homebrew.md](docs/homebrew.md) for the release and formula workflow.
