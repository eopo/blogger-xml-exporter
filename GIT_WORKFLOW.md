# Git Workflow & Branching Strategy

## 🎯 Strategy: GitHub Flow (Simplified)

This project uses **GitHub Flow** - optimized for continuous delivery on small to medium teams.

```
main (always deployable)
  ↑
  └─ feature branch (created from main)
     │
     ├─ Conventional commits
     ├─ PR with build verification
     └─ Merge back to main (triggers release)
```

### Why GitHub Flow?
- ✅ Simple: one main branch, feature branches only
- ✅ CI/CD friendly: PR checks before merge
- ✅ Release automation: merge to main = automatic versioning
- ✅ No complex branch management (no develop/release branches)

---

## 📋 Branch Naming

**Format**: `<type>/<short-description>`

```
feature/dark-mode           # New feature
fix/xml-encoding-bug        # Bug fix
docs/update-readme          # Documentation
chore/update-deps           # Dependencies/maintenance
```

---

## 🔄 Workflow: From Feature to Production

### 1. Create Feature Branch
```bash
git checkout main
git pull origin main
git checkout -b feature/my-feature
```

### 2. Conventional Commits (on feature branch)
```bash
git commit -m "feat(frontend): add dark mode toggle"
git commit -m "fix(backend): handle XML entities"
git push origin feature/my-feature
```

### 3. Open PR on GitHub
- ✅ Title: Clear description
- ✅ Labels: Set appropriate label (feature, bug, enhancement, etc.)
- ✅ Description: Why this change? (link issues with `Closes #123`)
- **Automatic**: `check.yml` runs (lint, test, build all components)
  - Wait for ✅ all checks pass

### 4. Review & Approval
- Code review from teammate
- Request changes if needed → push fixes (same branch)
- Checks re-run automatically

### 5. Merge to main
- **Merge strategy**: **Squash & Merge** (see below)
- Delete feature branch after merge

### 6. Automatic Release (release.yml triggers)
- Analyzes commits since last tag
- Bumps version (SemVer from commit types)
- Creates `v1.1.0` tag
- Generates GitHub Release (with PR links)
- Triggers `build.yml` → Docker image pushed

---

## 🔀 Merge Strategy: Squash & Merge

### Why Squash & Merge?

| Strategy                | History          | Use Case                         |
|-------------------------|------------------|----------------------------------|
| **Squash**              | Clean, linear    | ✅ Our approach (feature PRs)     |
| **Create Merge Commit** | Preserves branch | ❌ Complex multi-branch workflows |
| **Rebase**              | Rebased commits  | ❌ Less clear PR/commit tracking  |

### How to Squash Merge (GitHub UI)
1. Open PR on GitHub
2. Click "Squash and merge"
3. Commit message = PR title (auto-filled)
4. Click "Confirm squash and merge"

**Result**: 
```
main commit: "feat(frontend): add dark mode toggle" [from 5 feature commits]
```

Clean, single commit per feature = easy to read history, easy to revert if needed.

### Local Alternative (if needed)
```bash
git checkout main
git pull origin main
git merge --squash feature/my-feature
git commit -m "feat(frontend): add dark mode toggle"
git push origin main
```

---

## ⚡ Quick Operations

### Keep feature branch in sync with main
```bash
# During PR review, if main moved ahead
git fetch origin
git rebase origin/main  # or merge if you prefer
git push origin feature/my-feature --force-with-lease
```

### Local Git Config (Recommended)
```bash
# Set default pull to rebase (cleaner for feature branches)
git config pull.rebase true

# Set autosetuprebase for new branches tracking main
git config branch.autosetuprebase true

# Make it global (optional)
git config --global pull.rebase true
```

---

## ✅ Pre-Merge Build Verification

**Automatic on every PR:**

1. ✅ **Lint** (frontend + backend) - eslint, golangci-lint
2. ✅ **Test** (backend unit tests) - go test -race
3. ✅ **Build** (frontend bundle + Go binary)
4. ✅ **Coverage** (backend) - codecov upload

**In GitHub UI**: Shows as required check - **merge blocked** if any check fails.

### What You See:
```
✅ check / lint-backend
✅ check / lint-frontend
✅ check / test-backend
✅ check / build-frontend
✅ check / check-complete

→ Ready to merge!
```

---

## 🚫 Feature Branch Images: NOT Recommended

### Why NOT?
- ❌ Wastes registry space (every PR branch = new image)
- ❌ Manual cleanup needed (orphaned images)
- ❌ Adds cost (registry quota, build time)
- ❌ Not useful for small team (PR checks already verify build)

### Better Alternative:
- ✅ `check.yml` already builds all components on PR
- ✅ Verify locally: `make build` or `docker build` on your machine
- ✅ If PR checks pass → build is verified
- ✅ Only push image on merge to main → single `v1.1.0` image

---

## 🧪 Container Testing Strategy

- ✅ Unit tests + build verification sufficient
- ✅ Security scanning via Trivy on main  
- ❌ No redundant smoke tests ("docker run" is tautological)

---

### Recommended Approach: Minimal but Effective

**Two-tier approach:**



---

## 📌 Best Practices (Enforced)

✅ **Always**:
- Create feature branch from latest `main`
- Use conventional commits
- Add PR labels for categorization
- Wait for check.yml to pass before merge
- Use squash merge for clean history

❌ **Never**:
- Commit directly to `main` (always PR)
- Force-push to `main`
- Skip PR reviews
- Merge with failing checks
- Leave orphaned feature branches (delete after merge)

---

## 🔍 Troubleshooting

### "PR checks are failing, what do I do?"
1. Click on failing check → see error details
2. Fix locally: `git commit -m "fix: ..."`
3. Push same branch: `git push origin feature/my-feature`
4. Checks re-run automatically

### "I need to update main while my PR is open"
```bash
git fetch origin
git rebase origin/main
git push origin feature/my-feature --force-with-lease
```
This is safe (force-with-lease prevents accidents).

### "Accidental commit to main?"
```bash
git reset --soft HEAD~1  # Undo last commit, keep changes
git stash  # Save changes
git checkout -b fix/something  # Create proper branch
git stash pop  # Apply changes
```

---

## 📊 Workflow Summary

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Create branch: git checkout -b feature/x                │
│ 2. Commit: git commit -m "feat(scope): description"         │
│ 3. Push: git push origin feature/x                          │
│ 4. Open PR on GitHub (sets label)                           │
│                                                              │
│    ⏳ Automatic checks run (check.yml):                       │
│       ✅ lint-backend, lint-frontend                         │
│       ✅ test-backend, build-frontend                        │
│       ✅ check-complete gate                                 │
│                                                              │
│ 5. Code review + approval                                   │
│ 6. Squash & Merge → main                                    │
│                                                              │
│    ⏳ Automatic release (release.yml):                        │
│       ✅ Semantic versioning (feat → v1.1.0)                │
│       ✅ GitHub Release created with PR links                │
│       ✅ Docker build & push (build.yml)                    │
│       ✅ Container startup test                             │
│                                                              │
│ 7. Done! Production on v1.1.0 🚀                            │
└─────────────────────────────────────────────────────────────┘
```

This is **GitHub Flow**: simple, effective, and automated. ✨
