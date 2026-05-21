# Fix Auto-Tagging Pipeline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix missing tags v0.1.3/v0.1.4 and replace release-please with a simple auto-bump-patch workflow triggered on every push to main.

**Architecture:** Remove release-please (only fires on `fix:`/`feat:` commits). Replace with a GitHub Actions workflow that increments the patch version on every push to main, regardless of commit type. The existing `release.yaml` (goreleaser on `v*` tags) stays unchanged.

**Tech Stack:** GitHub Actions, bash, git tags, goreleaser (unchanged)

---

## Root Cause Summary

`release-please` only creates release PRs for `fix:`, `feat:`, and breaking changes. PRs #3 and #4 used `chore:` prefix. Result: no release PR created, no tag, manifest stuck at `0.1.2`.

The fix replaces the gating mechanism entirely. Every push to `main` = patch bump.

---

## File Map

| Action | File |
|--------|------|
| Delete | `.github/workflows/release-please.yaml` |
| Delete | `release-please-config.json` |
| Delete | `.release-please-manifest.json` |
| Create | `.github/workflows/auto-tag.yaml` |

`release.yaml` and `checks.yaml` are untouched.

---

### Task 1: Verify current state before making changes

**Files:**
- Read: `.release-please-manifest.json`
- Read: `.github/workflows/release-please.yaml`

- [ ] **Step 1: Confirm manifest version and latest tag**

```bash
cat .release-please-manifest.json
git tag --sort=-version:refname | head -5
```

Expected: manifest shows `"0.1.2"`, latest tag is `v0.1.2`.

- [ ] **Step 2: Confirm no unreleased tags exist**

```bash
git log --oneline --decorate | head -10
```

Expected: no `v0.1.3` or `v0.1.4` tags in output.

---

### Task 2: Create tag v0.1.4 manually

Tags are pushed to the remote — confirm with user before running `git push`.

- [ ] **Step 1: Create tag v0.1.4 on current HEAD**

```bash
git tag v0.1.4
```

Expected: no output (success).

- [ ] **Step 2: Verify tag exists locally**

```bash
git tag --sort=-version:refname | head -5
```

Expected: `v0.1.4` at top of list.

- [ ] **Step 3: Ask user to confirm push, then push tag**

```bash
git push origin v0.1.4
```

Expected: goreleaser release workflow triggers in GitHub Actions for `v0.1.4`.

---

### Task 3: Replace release-please with auto-bump-patch workflow

- [ ] **Step 1: Remove release-please files**

Delete these three files:
- `.github/workflows/release-please.yaml`
- `release-please-config.json`
- `.release-please-manifest.json`

```bash
rm .github/workflows/release-please.yaml release-please-config.json .release-please-manifest.json
```

Expected: files gone, `git status` shows 3 deletions.

- [ ] **Step 2: Create auto-tag workflow**

Create `.github/workflows/auto-tag.yaml`:

```yaml
name: Auto Tag

on:
  push:
    branches:
      - main

permissions:
  contents: write

jobs:
  auto-tag:
    name: Bump patch version and tag
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Skip if HEAD already tagged
        id: check
        run: |
          if git describe --exact-match HEAD 2>/dev/null; then
            echo "skip=true" >> "$GITHUB_OUTPUT"
          else
            echo "skip=false" >> "$GITHUB_OUTPUT"
          fi

      - name: Compute next patch tag
        if: steps.check.outputs.skip == 'false'
        id: bump
        run: |
          latest=$(git tag --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -1)
          if [ -z "$latest" ]; then
            latest="v0.0.0"
          fi
          major=$(echo "$latest" | cut -d. -f1 | tr -d 'v')
          minor=$(echo "$latest" | cut -d. -f2)
          patch=$(echo "$latest" | cut -d. -f3)
          new_tag="v${major}.${minor}.$((patch + 1))"
          echo "new_tag=$new_tag" >> "$GITHUB_OUTPUT"
          echo "Bumping from $latest to $new_tag"

      - name: Create and push tag
        if: steps.check.outputs.skip == 'false'
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git tag "${{ steps.bump.outputs.new_tag }}"
          git push origin "${{ steps.bump.outputs.new_tag }}"
```

- [ ] **Step 3: Verify no syntax errors (local lint)**

```bash
python3 -c "import yaml, sys; yaml.safe_load(open('.github/workflows/auto-tag.yaml'))" && echo "YAML valid"
```

Expected: `YAML valid`.

- [ ] **Step 4: Stage and commit changes**

```bash
git add .github/workflows/auto-tag.yaml
git add -u release-please-config.json .release-please-manifest.json .github/workflows/release-please.yaml
git status
```

Expected: 3 deletions + 1 new file staged.

- [ ] **Step 5: Ask user to confirm commit message, then commit**

```bash
git commit -m "chore: replace release-please with auto-bump-patch workflow

Every push to main now bumps the patch version automatically,
regardless of conventional commit type. Removes release-please
config files that only triggered on fix:/feat: commits."
```

---

### Task 4: Verify end-to-end behavior

- [ ] **Step 1: Confirm local tag state**

```bash
git tag --sort=-version:refname | head -5
```

Expected: `v0.1.4` is latest.

- [ ] **Step 2: Confirm workflow file exists**

```bash
cat .github/workflows/auto-tag.yaml
```

Expected: file present with correct content.

- [ ] **Step 3: Confirm release-please files removed**

```bash
ls release-please-config.json .release-please-manifest.json .github/workflows/release-please.yaml 2>&1
```

Expected: `No such file or directory` for all three.

- [ ] **Step 4: After user pushes the commit to main, verify in GitHub Actions**

In GitHub UI, confirm:
1. `auto-tag.yaml` workflow runs on the push
2. It creates tag `v0.1.5`
3. `release.yaml` goreleaser workflow triggers on `v0.1.5`

---

## What changes

| Before | After |
|--------|-------|
| release-please creates release PRs only for `fix:`/`feat:` | Every push to `main` bumps patch |
| Manifest file tracks version | Tags in git are the source of truth |
| 3 release-please config files | 1 auto-tag workflow file |
| `checks.yaml` and `release.yaml` | Unchanged |

## What is NOT verified locally

- GitHub Actions execution (requires push to remote)
- goreleaser trigger on v0.1.5 (depends on remote run)
- HOMEBREW_TAP_GITHUB_TOKEN secret still present in repo settings (unchanged, not our concern)
