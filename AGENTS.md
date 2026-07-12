# AGENTS.md — Project Overview & Guidelines

This document is for agents and developers working on this project. Start here.

## 📊 Project Structure

**Monorepo**: Single repository, two components

```
blogger-xml-exporter/
├── backend/              # Go 1.26 HTTP server
│   ├── main.go
│   ├── go.mod
│   └── internal/
│       ├── blogger/      # Blogger API client
│       ├── config/       # YAML config parsing
│       ├── httpapi/      # HTTP handlers
│       └── xmlgen/       # XML generation
├── frontend/             # Vue 3 + TypeScript
│   ├── src/
│   │   ├── App.vue
│   │   ├── components/   # Form components
│   │   └── composables/
│   ├── package.json
│   ├── vite.config.ts
│   └── vitest.config.ts
├── scripts/              # Utility scripts
├── .github/workflows/    # CI/CD pipelines
├── Makefile              # Development orchestration
├── Dockerfile            # Multi-stage Docker build
├── .version              # SemVer single source of truth
├── .commitlintrc.json    # Conventional Commits rules
└── README.md
```

## 🔨 Makefile & Tool Management

### Key Pattern: Reproducible Tool Installation

Tools are installed to `.bin/` directory (Kubernetes/Temporal pattern):

```makefile
AIR_VERSION := v1.65.3
GOLANGCI_LINT_VERSION := v2.12.2

go-install-tool = @command -v $(1) > /dev/null || \
  GOBIN="$(PWD)/.bin" go install $(2)@$(3)

setup:
  $(call go-install-tool,air,github.com/cosmtrek/air,$(AIR_VERSION))
  $(call go-install-tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))
  npm --prefix frontend ci
```

### Main Targets

| Target              | Purpose                           | Example                                    |
|---------------------|-----------------------------------|--------------------------------------------|
| `make setup`        | Install all tools + dependencies  | Run once after clone                       |
| `make dev`          | Start both servers (parallel)     | `make dev` → backend :8080, frontend :5173 |
| `make lint`         | Lint all code (Go + Vue)          | Runs before commit                         |
| `make test-backend` | Go unit tests with coverage       | `go test -race -coverprofile=coverage.out` |
| `make build`        | Build Go binary with version info | Static binary with ldflags                 |
| `make build-docker` | Multi-arch Docker image           | Builds for linux/amd64 + linux/arm64       |
| `make clean`        | Remove build artifacts            | Cleans .bin/, dist/, etc                   |

### Version Injection

```bash
make build
# Reads from .version, injects:
#   -X main.Version=1.0.0
#   -X main.CommitSHA=abc123
#   -X main.BuildTime=2026-07-12T...
```

## 📌 Versioning System

### Single Source of Truth: `.version`

```
1.0.0
```

Format: `MAJOR.MINOR.PATCH` (SemVer)

### Automatic Versioning

- `release.yml` analyzes **Conventional Commits**
- Calculates next version (SemVer bump)
- Updates `.version` file
- Creates Git tag + GitHub Release
- **No CHANGELOG.md file** — GitHub Releases are the only changelog

### Version Injection Points

1. **Go Binary** (via Makefile ldflags)
   ```go
   var Version = "dev"  // Injected at build time
   ```

2. **Docker Image** (via build args)
   ```dockerfile
   ARG VERSION=dev
   RUN go build -ldflags "-X main.Version=${VERSION}" ...
   ```

3. **Docker Tags**
   - `ghcr.io/repo:v1.0.0` (semantic tag)
   - `ghcr.io/repo:latest` (on main branch)

## 🌳 Git Workflow (GitHub Flow)

### Branching Strategy

```
main (production-ready)
  ↑
  └─ feature/* (feature branches from main)
     └─ squash-merge back to main
```

### Branch Naming

- `feature/description` — New features
- `fix/description` — Bug fixes  
- `docs/description` — Documentation
- `chore/description` — Build/deps/CI changes

### Workflow Steps

1. **Create branch**: `git checkout -b feature/my-feature`
2. **Make commits**: Use **Conventional Commits** format
   ```bash
   git commit -m "feat(backend): add health endpoint"
   git commit -m "fix(frontend): resolve form validation bug"
   ```
3. **Open PR** on GitHub
4. **Pre-merge checks** (check.yml runs):
   - ✅ Lint (ESLint, golangci-lint)
   - ✅ Unit tests (Go -race flag)
   - ✅ Build (frontend + backend)
5. **Merge to main**: Squash merge (clean linear history)
6. **Automatic release** (release.yml triggers):
   - Version bump (SemVer)
   - Git tag + GitHub Release
   - Docker build (build.yml triggers)

### No Feature Branch Images

❌ DON'T build Docker images per PR — waste of resources

✅ DO: Let PR checks verify build works, only push image on main merge

## 📋 Conventional Commits

### Format

```
<type>(<scope>): <description>
```

### Valid Types (affects version bump)

| Type                             | Bump  | Example                                   |
|----------------------------------|-------|-------------------------------------------|
| `feat`                           | MINOR | "feat(api): add export scheduling"        |
| `fix`                            | PATCH | "fix(parser): handle UTF-8 correctly"     |
| `feat!` / `fix!`                 | MAJOR | "feat(api)!: rename /export to /generate" |
| `docs`, `style`, `test`, `chore` | NONE  | No version bump                           |

### Examples

```bash
# Feature (MINOR bump)
git commit -m "feat(frontend): add dark mode toggle"

# Bug fix (PATCH bump)
git commit -m "fix(backend): resolve XML parsing error"

# Breaking change (MAJOR bump)
git commit -m "feat(api)!: remove deprecated endpoints"

# With body
git commit -m "feat(export): add --dry-run flag

Allows users to preview export without generating files.

Closes #123"
```

## 🔄 CI/CD Workflows

### check.yml (Pre-Merge Gate)

Runs on: Every push + PR to main

```
lint-backend ──┐
lint-frontend ─┼─→ check-complete (gate)
test-backend   │
test-frontend  │
build-frontend ┘
```

**Must pass before merge is allowed**

### release.yml (Semantic Versioning)

Runs on: Push to main

1. Analyzes Conventional Commits since last tag
2. Calculates version bump (SemVer)
3. Updates `.version` file
4. Creates Git tag + GitHub Release (auto-generated notes)
5. Triggers build.yml

### build.yml (Docker Build & Push)

Runs on: After release.yml succeeds

1. Builds multi-arch Docker image (amd64 + arm64)
2. Scans with Trivy (security)
3. Pushes to ghcr.io

### security.yml (Security Scans)

Runs on: Daily + on push

- `govulncheck` (Go CVE scanning)
- `npm audit` (JavaScript dependencies)
- `CodeQL` (code analysis)

### commitlint.yml (Commit Validation)

Runs on: PR + push to main

Enforces **Conventional Commits** format

## 🎯 Core Rules (From Project Guidelines)

### 1. Best Practices ALWAYS

❌ DON'T: Preserve existing patterns if they violate best practices

✅ DO: Follow best practices (DRY, KISS, YAGNI, 12-Factor)

### 2. Analyze Before Implementation

❌ DON'T: Blindly implement without critical review

✅ DO: Verify requirements, check if solution is necessary, document reasoning

**Example**: Don't add container tests just because they exist elsewhere — verify they add value

### 3. Pragmatic for Project Size

❌ DON'T: Over-engineer (full integration test suites, complex abstractions)

✅ DO: Size solutions to project complexity

**Current decision**: Unit tests + linting sufficient; no redundant smoke tests

### 4. Avoid Tautological Tests

❌ DON'T: `docker run --help` (redundant with docker build succeeding)

✅ DO: Test actual functionality or skip test

### 5. No Feature Branch Images

❌ DON'T: Build Docker image per PR branch

✅ DO: Only build on main merge (resource efficiency)

### 6. Version Pinning: Major.Minor Only

All GitHub Actions use major.minor versioning (`@v7`, `@v6`, `@v9`):

- ✅ Auto-get patch updates
- ✅ Auto-get minor updates
- ✅ Stability (no breaking major changes)

Example:
```yaml
uses: actions/checkout@v7          # NOT v7.0.0
uses: actions/setup-go@v6          # NOT v6.5.0
uses: golangci/golangci-lint-action@v9  # NOT v9.3.0
```

## 🏗️ Architecture

### Backend (Go 1.26)

- **HTTP Server** on :8080
- **Static binary** (CGO_ENABLED=0, distroless container)
- **Packages**:
  - `blogger/` — Blogger API client
  - `config/` — YAML configuration
  - `httpapi/` — HTTP route handlers
  - `xmlgen/` — XML file generation

### Frontend (Vue 3 + TypeScript)

- **Dev Server** on :5173 (Vite + HMR)
- **Proxies** to backend :8080
- **Tailwind CSS** v4 (PostCSS)
- **Components**: Form fields, arrays, combobox (Tom Select)
- **Build**: `npm run build` → dist/ (optimized)

### Docker

- **Stage 1**: Node 22 Alpine (frontend build)
- **Stage 2**: Go 1.24 Alpine (backend build)
- **Final**: `gcr.io/distroless/static-debian12:nonroot` (security-hardened, no shell)

## 📝 Code Guidelines

### Comments (Per instructions)

✅ DO:
- JSDoc/DOCSTRING for functions, classes
- Explain WHY, not HOW
- Keep language consistent (same language as rest of project)

❌ DON'T:
- Comment every line
- Explain implementation details
- Use different language than codebase

### Error Handling (Go)

✅ DO: Explicit error handling
```go
if err != nil {
    return err  // or handle appropriately
}
```

❌ DON'T: Ignore errors
```go
_ = resp.Body.Close()  // Tautological, doesn't handle error
```

Proper:
```go
defer func() {
    if err := resp.Body.Close(); err != nil {
        log.Printf("error closing response: %v", err)
    }
}()
```

## 🚀 Common Tasks

### Adding a Feature

1. Create branch: `git checkout -b feature/my-feature`
2. Make changes + commit with conventional format
3. Open PR
4. Wait for check.yml to pass
5. Squash merge to main
6. release.yml auto-triggers (version + tag)
7. build.yml auto-triggers (Docker build)

### Updating Dependencies

**Backend (Go)**:
```bash
go get -u ./...
go mod tidy
```

**Frontend (Vue)**:
```bash
npm update --prefix frontend
npm audit fix --prefix frontend
```

Commit as: `chore: update dependencies`

### Debugging Locally

```bash
make dev                           # Start both servers
curl http://localhost:8080/...     # Test backend
# Open http://localhost:5173       # Test frontend
```

### Building Release

1. Merge PR to main
2. release.yml auto-runs:
   - Detects commits
   - Bumps version
   - Creates tag
   - Creates GitHub Release
3. build.yml auto-runs:
   - Builds Docker image
   - Pushes to registry

No manual steps needed!

## ⚠️ What NOT to Do

1. ❌ Manually edit `.version` file
2. ❌ Force-push to `main`
3. ❌ Create long-running feature branches (use small PRs)
4. ❌ Skip pre-merge checks (check.yml must pass)
5. ❌ Build Docker images per PR (only on main)
6. ❌ Add redundant tests that don't test real behavior
7. ❌ Hardcode configuration values
8. ❌ Ignore error handling

## 📚 References

- [GIT_WORKFLOW.md](GIT_WORKFLOW.md) — Detailed git workflow & branching
- [VERSIONING.md](VERSIONING.md) — Version management & release process
- [Makefile](Makefile) — Build configuration & tool management
- [README.md](README.md) — Project overview & quick start
