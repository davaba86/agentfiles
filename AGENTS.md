@.agent-config/AGENTS.md

# Agent Instructions

## Git workflow — override

Do not create branches, push commits, or open PRs unless user explicitly asks. Never push directly to `main` or `master` — even if user says "proceed" to a plan that includes it, stop and ask explicitly. Make changes locally only. Before any git command (branching, committing, pushing, PR creation): ask user for approval first.

## Planning files

Save implementation plans as `PLAN-<description>.md` at repo root.

## Phased projects

Split work into phases before starting. Complete one phase, then stop and ask user to confirm before next. If a `PLAN-*.md` exists, update it after each phase.
