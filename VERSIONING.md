# Versioning & Release Strategy

This project follows **Semantic Versioning (SemVer)** and uses **Conventional Commits** for automated version management.

## 🏷️ Version Management

### Current Version
- Version stored in: **`.version`** file (single source of truth)
- Format: `MAJOR.MINOR.PATCH` (e.g., `1.0.0`)
- Version injected into:
  - Go binary (ldflags: `main.Version`, `main.CommitSHA`, `main.BuildTime`)
  - Docker images (build args + image tags)
  - GitHub releases (automated)

### Version Sources
```bash
# Read current version
cat .version

# Get version from running binary
./bin/blogger-xml-exporter --version  # (if implemented)
```

## 📝 Changelog: GitHub Releases Only

This project uses **GitHub Releases** as the single changelog source:

- ✅ Auto-generated from PR titles & labels
- ✅ Links to all PRs merged in the release
- ✅ Automatic contributor attribution
- ✅ Accessible at: `https://github.com/owner/repo/releases`

**No CHANGELOG.md file** — it's redundant with auto-generated release notes.
- ✅ Interactive GitHub UI with full PR context
- 📍 Generated automatically by GitHub Actions

**Example:**
```
## What's Changed
✨ Features
- Add dark mode toggle by @alice (#45)
- Implement export scheduling by @bob (#46)

🐛 Bug Fixes
- Fix XML parsing for special characters by @charlie (#47)

👥 New Contributors
- @alice made their first contribution in #45



This project uses **Conventional Commits** to enable automatic semantic versioning.

### Commit Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Valid Types
- **feat**: A new feature → triggers `MINOR` bump
- **fix**: A bug fix → triggers `PATCH` bump
- **feat!**: Breaking change → triggers `MAJOR` bump
- **fix!**: Breaking change → triggers `MAJOR` bump
- **docs**: Documentation only
- **style**: Code style (no functional changes)
- **refactor**: Code refactor (no functionality change)
- **perf**: Performance improvement
- **test**: Test additions/modifications
- **chore**: Build/dependency changes
- **ci**: CI/CD changes
- **revert**: Revert a previous commit

### Examples
```bash
# Feature (MINOR bump)
git commit -m "feat(frontend): add dark mode toggle"

# Bug fix (PATCH bump)
git commit -m "fix(backend): resolve XML parsing error on special chars"

# Breaking change (MAJOR bump)
git commit -m "feat(api)!: rename /export to /generate endpoint"

# With body/footer
git commit -m "feat(cli): add --dry-run flag

Allows users to preview export without generating files.

Closes #123"
```

## 🔄 Release Workflow

### Automatic Release (Recommended)
1. **Create PR** with conventional commits
2. **Merge to `main`** → `release.yml` workflow runs automatically
3. **Workflow actions**:
   - Analyzes commits since last tag
   - Calculates next version (SemVer)
   - Updates `.version` file
   - Creates Git tag and GitHub release
   - Triggers `build.yml` for Docker build + push

### Manual Release (if needed)
```bash
# Trigger release workflow manually
gh workflow run release.yml --ref main --field bump=patch

# Or use semantic-release CLI locally
npx semantic-release --no-ci
```

## 📦 Build & Deployment

### Local Build with Version
```bash
make build  # Reads from .version, injects via ldflags
```

### Docker Build with Version
```bash
make build-docker  # Tags image as :latest and :v1.0.0
```

### GitHub Actions CI/CD
- **check.yml** → Lint & test (runs on every push/PR)
- **release.yml** → Semantic versioning & tag creation (runs on main merge)
- **build.yml** → Docker build & push (triggered after release)
- **security.yml** → Security scanning (daily + on push)
- **commitlint.yml** → Commit message validation (PR + push)

## 🏷️ Version Injection

Version information is injected at build time:

### Go Binary (via ldflags)
```go
// backend/main.go
var (
    Version   = "dev"        // Set via -X main.Version=1.0.0
    CommitSHA = "unknown"    // Set via -X main.CommitSHA=abc123
    BuildTime = "unknown"    // Set via -X main.BuildTime=2024-07-12T...
)
```

### Docker (via build args)
```dockerfile
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME=unknown

RUN go build \
    -ldflags "-X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA} -X main.BuildTime=${BUILD_TIME}" \
    ...
```

## 📊 Version Bumping Logic

```
Commit type       Current    Next        Reason
─────────────────────────────────────────────────────
fix               1.0.0  →   1.0.1       Patch (bug fix)
feat              1.0.0  →   1.1.0       Minor (new feature)
feat! / fix!      1.0.0  →   2.0.0       Major (breaking change)
docs/style/test   1.0.0  →   1.0.0       No change (dev commits)
```

## 🎯 Release Process (Automatic Example)

### Step 1: Create PR with conventional commits
```bash
git checkout -b feat/my-feature
git commit -m "feat(frontend): add new dashboard widget"
git commit -m "fix(backend): improve error handling"
git push origin feat/my-feature
```

### Step 2: Merge to main
```bash
# PR approved & merged to main
git checkout main && git pull
```

### Step 3: Automatic release triggered
The `release.yml` workflow:
1. Detects 1 feat + 1 fix commit
2. Calculates version: `1.0.0` → `1.1.0` (MINOR bump)
3. Creates Git tag `v1.1.0`
4. Creates GitHub Release with auto-generated notes:
   - 🔗 Links all PRs merged since last release
   - 👥 Lists contributors
   - 🏷️ Categorizes changes (Features, Bug Fixes, Chores)
5. Triggers `build.yml`
6. Builds & pushes Docker image as `ghcr.io/repo:v1.1.0` and `ghcr.io/repo:latest`

### GitHub Release (Auto-Generated)

**Release Body**:
```
## What's Changed
- Added new dashboard widget (#45) by @alice
- Improved error handling (#46) by @bob

## New Contributors
- @alice made their first contribution in #45
```

## 📌 Best Practices

✅ **DO:**
- Use conventional commit format in every commit
- **Add appropriate GitHub PR labels** for better release notes categorization
- Group related changes in single commit
- Keep commits atomic and logically separated
- Write descriptive commit bodies for complex changes
- Reference issues in commit footers: `Closes #123`

❌ **DON'T:**
- Mix unrelated changes in one commit
- Use vague commit messages
- Force-push to `main` (use squash-merge for PRs)
- Manually edit `.version` file
- Create releases without merging to `main` first
- **Forget to add PR labels** (makes release notes less organized)

## 🏷️ GitHub Release Configuration

GitHub's **Auto-Generate Release Notes** uses PR labels to categorize changes in the release:

### Default Categories (Auto-Recognized)
GitHub automatically recognizes standard labels:

| Label                   | Category            | Example                             |
|-------------------------|---------------------|-------------------------------------|
| `bug`, `bugfix`         | 🐛 Bug Fixes        | Fixed critical authentication issue |
| `feature`, `feat`       | ✨ New Features      | Added dark mode support             |
| `enhancement`           | 🚀 Enhancements     | Improved search performance         |
| `documentation`, `docs` | 📚 Documentation    | Updated API documentation           |
| `breaking-change`       | ⚠️ Breaking Changes | Removed deprecated endpoints        |
| `security`              | 🔒 Security         | Patched XSS vulnerability           |
| `chore`, `maintenance`  | 🔧 Maintenance      | Updated dependencies                |

### How to Use

**When creating a PR**, add appropriate label(s):
1. Open PR on GitHub
2. Add label(s) matching your change type
3. When PR is merged, GitHub's auto-generate will categorize it

**Example:**
```
PR Title: Add dark mode toggle
Labels: enhancement, feature
↓
Release Notes:
✨ New Features
- Add dark mode toggle (#123)
```

### Customizing Categories

To customize labels, create `.github/release.yml`:
```yaml
changelog:
  exclude:
    labels: [skip-changelog]
  categories:
    - title: 🎉 Highlights
      labels: [highlight]
    - title: 🐛 Bug Fixes
      labels: [bug, bugfix]
    - title: ✨ Features
      labels: [feature, feat, enhancement]
    - title: 📚 Documentation
      labels: [documentation, docs]
```

Currently we use **GitHub's default categories** (no custom config needed).
