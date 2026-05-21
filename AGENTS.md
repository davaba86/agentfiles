# Agent Instructions

## First files to read

Before changing code, read:

1. README.md
2. Taskfile.yml
3. Relevant docs under docs/
4. Relevant source files for the task

## Tooling

- Prefer project-defined task commands.
- Do not invent commands; inspect the repo first.

## Instruction files

- When updating AI instructions, always edit `AGENTS.md`.
- Never edit `CLAUDE.md`, `GEMINI.md`, or any other agent-specific instruction file. `AGENTS.md` is the sole source of truth for shared rules. Edits to other instruction files are forbidden unless the user explicitly asks.

## Planning files

- When the user asks to save the current chat or discussion into a PLAN file, save it in the repository root as `PLAN-SOMETHING-RELATED-TO-CHAT.md`.

## Code style

- Do not use emojis inside code unless explicitly asked to by the user.
- If tempted to use emojis, use text replacements such as `[ok]`, `[fail]`, or another clear ASCII marker.
- Do not use em dashes under any circumstances.
- Always use `.yaml` for YAML files, not `.yml`.

## Workflow

Before editing:
- Summarize the relevant project structure.
- State assumptions.
- Before removing or disabling existing functionality, such as CI steps, configuration, or integrations, confirm with the user. It may be intentional. Prefer asking over assuming it is dead code or a bug.

After editing:
- Run the smallest relevant validation.
- Summarize what changed.
- Mention anything not verified.

When working in a repository:
- Do not create branches, push commits, or open PRs unless the user explicitly asks.
- Never push directly to `main` or `master` under any circumstances unless the user explicitly asks.
- Make code changes locally only.
- When ready to run any git command (branching, committing, pushing, PR creation), ask the user for approval first and wait for confirmation before executing.
