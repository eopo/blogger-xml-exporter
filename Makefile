# Development workflow for Blogger XML Exporter
# Follows Go best practices: tools in .bin/ (Temporal/Kubernetes pattern)
# Run `make help` to list all available targets.

.DEFAULT_GOAL := help
.PHONY: help setup dev watch-backend watch-frontend stop build build-docker lint lint-fix format test clean

# Variables
ROOT := $(shell git rev-parse --show-toplevel)
LOCALBIN := $(ROOT)/.bin
BLOGGER_API_KEY ?= dev-dummy-key
CONFIG_PATH ?= ./config.yaml
DOCKER_IMAGE ?= blogger-xml-exporter:local
BIN_NAME := blogger-xml-exporter

# Add .bin to PATH so installed tools are found first
export PATH := $(LOCALBIN):$(PATH)

# Tool versions (pinned for reproducibility)
AIR_VERSION := v1.65.3
GOLANGCI_LINT_VERSION := v1.64.8

# Tool binaries
AIR := $(LOCALBIN)/air
GOLANGCI_LINT := $(LOCALBIN)/golangci-lint

help: ## Show this help message
	@grep -E '^[a-zA-Z_:-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

# Install tool if not already present (macro from Temporal/Kubernetes pattern)
define go-install-tool
@[ -f $(1) ] || { \
	set -e; \
	mkdir -p $(LOCALBIN); \
	package=$(2)@$(3); \
	tmpdir=$$(mktemp -d); \
	GOBIN=$${tmpdir} go install $${package}; \
	mv $${tmpdir}/$$(basename "$(1)") $(1); \
	rm -rf $${tmpdir}; \
}
endef

$(AIR): $(LOCALBIN)
	$(call go-install-tool,$(AIR),github.com/air-verse/air,$(AIR_VERSION))

$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

setup: $(AIR) $(GOLANGCI_LINT) ## Install npm + dev tools
	npm install
	@echo "✓ Setup complete: tools installed in .bin/"

lint-backend: $(GOLANGCI_LINT) ## Lint backend Go code
	@cd backend && $(GOLANGCI_LINT) run ./...

lint-frontend: ## Lint frontend (Vue/TypeScript)
	@cd frontend && npm run lint

lint: lint-frontend lint-backend ## Lint all (frontend + backend)

test-backend: ## Run backend tests
	@cd backend && go test -v ./...

test: test-backend ## Run all tests (backend only)

dev: ## Start Backend (watch) + Frontend (watch) together
	@echo "Starting Backend (watch) + Frontend (watch)..."
	@echo "Press Ctrl+C to stop all processes"
	@echo ""
	@make -j watch-backend watch-frontend

watch-backend: $(AIR) ## Start Backend with file watching (air)
	@cd backend && \
		BLOGGER_API_KEY=$(BLOGGER_API_KEY) \
		CONFIG_PATH=../$(CONFIG_PATH) \
		air

watch-frontend: ## Start Frontend with HMR (Vite)
	@cd frontend && npm run dev

stop: ## Stop all dev processes (air + vite)
	@pkill -f "air|vite" || echo "No processes to stop"
	@echo "✓ All dev processes stopped"

build: ## Build frontend (Vue 3 → web/static/) and Go binary
	npm run build
	cd backend && go build -o ../bin/$(BIN_NAME) .
	@echo "✓ Build complete (bin/$(BIN_NAME) + web/static/)"

build-docker: ## Build Docker image (complete multi-stage build)
	docker build -t $(DOCKER_IMAGE) .
	@echo "✓ Docker image built: $(DOCKER_IMAGE)"

lint-fix: ## Fix linting issues (ESLint + gofmt)
	npm run lint:fix
	@cd backend && go fmt ./...

format: ## Format code (Prettier + gofmt)
	npm run format
	@cd backend && go fmt ./...

clean: ## Clean build artifacts and tools
	rm -rf bin/ web/static/ coverage/ dist/ $(LOCALBIN)
	@cd backend && rm -rf tmp/
	@echo "✓ Clean complete"

