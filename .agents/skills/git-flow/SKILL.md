---
name: git-flow
description: A structured methodology for managing version control, remotes, and pull requests within the ZTVS repository.
---

# Skill: Git Flow for ZTVS

This skill provides a structured methodology for managing version control, remotes, and pull requests within the Zero Trust Vulnerability Scanner (ZTVS) repository.

## Overview
ZTVS standardizes on GitHub-based collaboration. The primary delivery mechanism is the Pull Request (PR).

## 🛠️ Step-by-Step Implementation

### 1. Branching Strategy
Always create a fresh branch for new work. 
- Use the conventional naming: `feat/<name>`, `fix/<name>`, or `docs/<name>`.
- Command: `git checkout -b feat/my-new-feature`

### 2. Committing Work
Follow the [Conventional Commits](.agents/rules/git.md) rules.
- Workflow: Use the `/git` workflow to stage and commit.
- Command: `git commit -m "feat(plugin): add network check"`

### 3. Pushing & Remotes
- Workflow: Use the `/git` workflow to push to your target remote.
- Identifying remotes: `git remote -v`.
- Important: Code MUST be pushed to a remote branch before a Pull Request can be created.

### 4. PR Creation via `gh`
Automate the pull request lifecycle using the GitHub CLI once your code is pushed:
- Workflow: Use the `/pr` workflow to create and view Pull Requests.
- Command: `gh pr create --base main --head <current-branch> --title "feat: my change" --body "Detailed description"`
- Link to RFCs if applicable: `Depends on RFC-001`.

### 5. Keeping Main Up-to-Date
- Checkout main: `git checkout main`
- Pull latest: `git pull <remote> main`
- Rebase your feature branch: `git checkout feat/my-feature && git rebase main`

## 🛡️ Best Practices
- **Atomic Commits**: Keep each commit focused on one logical change.
- **DCO/Signing**: If required, use `git commit -s` for Developer Certificate of Origin.
- **Squash-on-Merge**: Prefer squashing commits into a single descriptive commit upon PR merge to keep the `main` history linear and clean.
