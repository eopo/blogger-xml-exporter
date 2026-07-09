# Development workflow for Blogger XML Exporter.
# Run `make help` to list all available targets.

.DEFAULT_GOAL := help

.PHONY: help setup check-go install-air install-tailwind install-pre-commit dev \
	build build-go build-css build-image test test-coverage lint fmt pre-commit-install pre-commit-run

TAILWIND_BIN := tools/tailwindcss
GOBIN := $(shell go env GOBIN 2>/dev/null)
GOPATH := $(shell go env GOPATH 2>/dev/null)
AIR_BIN := $(if $(GOBIN),$(GOBIN)/air,$(GOPATH)/bin/air)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: check-go install-air install-tailwind install-pre-commit ## Setup development environment

check-go: ## Verify Go is installed
	@command -v go >/dev/null 2>&1 || { echo "Go not found, see https://go.dev/doc/install"; exit 1; }
	@echo "Go found: $$(go version)"

install-air: ## Install air (live reload for Go)
	@if [ -x "$(AIR_BIN)" ]; then \
		echo "air already installed: $(AIR_BIN)"; \
	else \
		go install github.com/air-verse/air@latest; \
	fi

install-tailwind: ## Download Tailwind CSS CLI
	@if [ -x "$(TAILWIND_BIN)" ]; then \
		echo "Tailwind CLI already available: $(TAILWIND_BIN)"; \
		exit 0; \
	fi; \
	os="$$(uname -s)"; arch="$$(uname -m)"; \
	case "$$os" in \
		Darwin) platform="macos" ;; \
		Linux) platform="linux" ;; \
		*) platform="" ;; \
	esac; \
	case "$$arch" in \
		arm64|aarch64) cpu="arm64" ;; \
		x86_64|amd64) cpu="x64" ;; \
		*) cpu="" ;; \
	esac; \
	if [ -z "$$platform" ] || [ -z "$$cpu" ]; then \
		echo "Could not detect OS/architecture ($$os/$$arch)"; \
		echo "Download Tailwind CLI manually: https://tailwindcss.com/blog/standalone-cli"; \
		exit 1; \
	fi; \
	asset="tailwindcss-$${platform}-$${cpu}"; \
	url="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$${asset}"; \
	echo "Downloading Tailwind CLI (latest, $${asset})..."; \
	mkdir -p "$(dir $(TAILWIND_BIN))"; \
	if ! curl -fsSL -o "$(TAILWIND_BIN)" "$$url"; then \
		echo "Download failed. Install Tailwind CLI manually: https://tailwindcss.com/blog/standalone-cli"; \
		rm -f "$(TAILWIND_BIN)"; \
		exit 1; \
	fi; \
	chmod +x "$(TAILWIND_BIN)"; \
	echo "Tailwind CLI installed: $(TAILWIND_BIN)"

install-pre-commit: ## Install pre-commit hooks
	@command -v pre-commit >/dev/null 2>&1 || { \
		echo "pre-commit not found. Install with: pip install pre-commit"; \
		exit 1; \
	}
	pre-commit install
	pre-commit install --hook-type commit-msg
	@echo "Pre-commit hooks installed"

pre-commit-run: ## Run all pre-commit hooks
	pre-commit run --all-files

dev: ## Start dev server (air + Tailwind watch)
	@[ -x "$(AIR_BIN)" ] || { echo "air not found - run 'make setup' first"; exit 1; }
	@[ -x "$(TAILWIND_BIN)" ] || { echo "Tailwind CLI not found - run 'make setup' first"; exit 1; }
	@trap 'kill 0' EXIT; \
	$(TAILWIND_BIN) -i web/tailwind.src.css -o web/static/css/style.css --watch & \
	$(AIR_BIN)

build-go: ## Build Go binary
	go build -o bin/blogger-xml-exporter .

build-css: install-tailwind ## Compile Tailwind CSS (minified)
	$(TAILWIND_BIN) -i web/tailwind.src.css -o web/static/css/style.css --minify

build: build-go build-css ## Build Go binary and CSS

build-image: build-css ## Build Docker image
	docker build -t blogger-xml-exporter .

test: ## Run Go tests
	go test -v ./...

test-coverage: ## Run Go tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run Go linter (golangci-lint)
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	}
	golangci-lint run --timeout=5m

fmt: ## Format Go code
	goimports -w .
	gofmt -w .
