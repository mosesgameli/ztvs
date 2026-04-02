---
name: pr-flow
description: Create, view, and manage Pull Requests on GitHub.
---


1. Generate a Pull Request via GitHub CLI:
// turbo
`gh pr create --base main --head <current-branch> --title "<type>: <summary>" --body "Implemented according to roadmap."`

2. Verify PR status:
// turbo
`gh pr status`

3. List existing pull requests:
// turbo
`gh pr list`

4. View a specific pull request:
// turbo
`gh pr view <pr-number>`
