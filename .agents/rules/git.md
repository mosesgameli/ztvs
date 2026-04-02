 
trigger: always_on
 

# ZTVS Git Standards & Best Practices

These rules govern the Git lifecycle for the ZTVS project to ensure a clean commit history and automated CI/CD compatibility.

## 📡 Remotes & Pushing
1.  **Flexible Remotes**: You may use any remote (e.g., `me`, `origin`, `upstream`). 
2.  **Explicit Push**: Always specify the target remote and branch during push to avoid ambiguity (e.g., `git push <remote> <branch>`).
3.  **PR-First**: All code MUST be merged via Pull Requests. Direct pushes to `main` are discouraged.

## 📝 Commit Messages
We follow [Conventional Commits](https://www.conventionalcommits.org/).

Format: `<type>(<scope>): <description>`

Common types:
- `feat`: A new feature (e.g., `feat(engine): add parallel scan support`).
- `fix`: A bug fix (e.g., `fix(sdk): resolve handshake timeout`).
- `docs`: Documentation only changes (e.g., `docs(readme): add installation guide`).
- `refactor`: A code change that neither fixes a bug nor adds a feature.
- `chore`: Updating build tasks, package manager configs, etc.

## 🌿 Branching Strategy
Naming convention: `<type>/<summary>`

Examples:
- `feat/protocol-handshake`
- `fix/rpc-marshal-error`
- `docs/api-specification`

## 🏁 Pull Requests
1.  **Clear Title**: Use the conventional commit style for PR titles.
2.  **Detailed Body**: Describe *what* changed and *why*. Link to any relevant RFCs or issues.
3.  **Clean History**: Rebase and squash commits before merging if the history is cluttered with "fixup" commits.
