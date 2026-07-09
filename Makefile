# Development workflow for Blogger XML Exporter
# Run `make help` to list all available targets.

.DEFAULT_GOAL := help
.PHONY: help setup setup-css-tools dev build build-go build-css build-docker lint test test-coverage fmt clean


help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup-css-tools: ## Download Tailwind CSS CLI (auto-detects OS/arch/libc)
	@mkdir -p tools
	@OS=$$(uname -s); ARCH=$$(uname -m); \
	case "$$OS" in \
	  Darwin) os=macos ;; \
	  Linux)  os=linux ;; \
	  MINGW*|MSYS*|CYGWIN*) os=windows ;; \
	  *) echo "ERROR: unsupported OS '$$OS'"; exit 1 ;; \
	esac; \
	case "$$ARCH" in \
	  arm64|aarch64) arch=arm64 ;; \
	  x86_64|amd64)  arch=x64 ;; \
	  armv7l|armv6l) echo "INFO: 32-bit ARM not supported for prebuilt Tailwind. Skipping."; exit 0 ;; \
	  *) echo "ERROR: unsupported arch '$$ARCH'"; exit 1 ;; \
	esac; \
	libc=""; \
	if [ "$$os" = linux ] && [ -f /lib/ld-musl-$$ARCH.so.1 ]; then libc="-musl"; fi; \
	ext=""; out="tools/tailwindcss"; \
	if [ "$$os" = windows ]; then ext=".exe"; out="tools/tailwindcss.exe"; fi; \
	BINARY="tailwindcss-$$os-$$arch$$libc$$ext"; \
	curl -sSfL "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$$BINARY" -o "$$out"; \
	[ "$$os" = windows ] || chmod +x "$$out"; \
	echo "✓ Tailwind ready ($$BINARY)"

setup: ## Setup development environment
	@mkdir -p tools
	# Install dev tools locally in tools/
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b ./tools
	GOBIN=$(shell pwd)/tools go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./tools latest
	$(MAKE) setup-css-tools
	@echo "Setup complete. Run 'make dev' to start development server."

dev: ## Start development server with live reload
	@[ -x tools/air ] || (echo "Run 'make setup' first"; exit 1)
	@trap 'kill 0' EXIT; \
	tools/tailwindcss -i web/tailwind.src.css -o web/static/css/style.css --watch & \
	./tools/air

build: build-go build-css ## Build Go binary and CSS (local only, for development)

build-docker: ## Build Docker image locally (all compilation happens in container)
	docker build -t blogger-xml-exporter:local .
	@echo "Docker image built: blogger-xml-exporter:local"

build-go: ## Build Go binary
	@mkdir -p bin
	go build -o bin/blogger-xml-exporter .
	@echo "Binary built: bin/blogger-xml-exporter"

build-css: ## Compile Tailwind CSS (minified)
	@[ -x tools/tailwindcss ] || (echo "Tailwind CLI not found. Run 'make setup'"; exit 1)
	@tools/tailwindcss -i web/tailwind.src.css -o web/static/css/style.css --minify
	@echo "CSS compiled: web/static/css/style.css"

test: ## Run Go tests
	go test -v ./...

test-coverage: ## Run Go tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run Go linter (golangci-lint)
	@[ -x tools/golangci-lint ] || { echo "golangci-lint not found in ./tools. Run 'make setup' first"; exit 1; }
	./tools/golangci-lint run

fmt: ## Format Go code
	@[ -x tools/goimports ] || (echo "goimports not found. Run 'make setup' first"; exit 1)
	./tools/goimports -w .
	gofmt -w .
	@echo "Code formatted"

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html
	@echo "Clean complete"
