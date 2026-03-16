---
description: This workflow automates the process of branching, committing, and opening a PR using the GitHub CLI.
---

# Workflow: Ship Feature

This workflow automates the process of branching, committing, and opening a PR using the GitHub CLI.

## Steps

1. **Safety Check**: Ensure the agent is on the `main` branch and has the latest code.
   - Command: `git checkout main && git pull origin main`

2. **Branching**: Create a unique feature branch based on the user's task.
   - Command: `git checkout -b feat/{{branch_name}}`

3. **Implementation**: Apply the requested changes to the codebase (e.g., updating .gitignore).

4. **Verification**: Run a quick check (if applicable).
   - Command: `go fmt ./...` (Optional for Go projects)

5. **Commit & Push**:
   - Command: `git add .`
   - Command: `git commit -m "{{commit_message}}"`
   - Command: `git push origin feat/{{branch_name}}`

6. **GitHub PR**: Use the GitHub CLI to create the PR and provide the link.
   - Command: `gh pr create --title "{{commit_message}}" --body "Automated PR via Antigravity Agent: {{pr_description}}" --web`
   - Note: "The PR is now open. Please review and merge it manually in the browser."
