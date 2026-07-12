# Git Configuration for Developer Setup

Run these commands once to configure your local Git for optimal workflow:

```bash
# Use rebase for pulls (cleaner feature branch history)
git config pull.rebase true

# Auto-setup rebase for new branches
git config branch.autosetuprebase true

# Better merge conflict resolution
git config merge.conflictstyle diff3

# Safer force-push (use --force-with-lease instead of --force)
git config push.default simple

# Apply these globally (optional):
git config --global pull.rebase true
git config --global branch.autosetuprebase true
git config --global merge.conflictstyle diff3
```

## What These Do

| Setting | Effect |
|---------|--------|
| `pull.rebase true` | `git pull` = rebase instead of merge (cleaner history) |
| `branch.autosetuprebase true` | New branches track with rebase by default |
| `merge.conflictstyle diff3` | Better context in conflict markers |
| `push.default simple` | Safe: only push current branch |

## Result

```bash
# Your workflow stays the same, but cleaner:
git checkout -b feature/x
git commit ...
git push origin feature/x
git pull origin main  # = rebase (clean), not merge
```
