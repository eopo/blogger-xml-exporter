# Development workflow for Blogger XML Exporter
# Run `make help` to list all available targets.

.DEFAULT_GOAL := help
.PHONY: help setup dev build build-go build-css lint test test-coverage fmt clean

# Tool versions (following Go best practices)
TAILWIND_VERSION := latest
GOLANGCI_LINT_VERSION := latest

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup development environment
	go install github.com/air-verse/air@latest
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.12.2
	@mkdir -p tools
	@curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64 -o tools/tailwindcss || \
		curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tools/tailwindcss
	@chmod +x tools/tailwindcss
	@echo "Setup complete. Run 'make dev' to start development server."

dev: ## Start development server with live reload
	@[ -x tools/tailwindcss ] || (echo "Run 'make setup' first"; exit 1)
	@trap 'kill 0' EXIT; \
	tools/tailwindcss -i web/tailwind.src.css -o web/static/css/style.css --watch & \
	air

build: build-go build-css ## Build Go binary and CSS

build-go: ## Build Go binary
	@mkdir -p bin
	go build -o bin/blogger-xml-exporter .
	@echo "Binary built: bin/blogger-xml-exporter"

build-css: ## Compile Tailwind CSS (minified)
	@[ -x tools/tailwindcss ] || (echo "Tailwind CLI not found. Run 'make setup'"; exit 1)
	tools/tailwindcss -i web/tailwind.src.css -o web/static/css/style.css --minify
	@echo "CSS compiled: web/static/css/style.css"

test: ## Run Go tests
	go test -v ./...

test-coverage: ## Run Go tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run Go linter (golangci-lint)
	golangci-lint run --timeout=5m

fmt: ## Format Go code
	goimports -w .
	gofmt -w .
	@echo "Code formatted"

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html
	@echo "Clean complete"
