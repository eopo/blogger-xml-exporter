# Development workflow for Blogger XML Exporter
# Follows Gitea pattern: make dev for watch mode, make stop to cleanup
# Run `make help` to list all available targets.

.DEFAULT_GOAL := help
.PHONY: help setup dev watch-backend watch-frontend stop build build-docker lint lint-fix format test clean

# Variables
BLOGGER_API_KEY ?= dev-dummy-key
CONFIG_PATH ?= ./config.yaml
DOCKER_IMAGE ?= blogger-xml-exporter:local
BIN_NAME := blogger-xml-exporter

help: ## Show this help message
	@grep -E '^[a-zA-Z_:-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Install npm + Go dependencies
	npm install
	cd backend && go mod download
	@echo "✓ Dependencies installed"

dev: ## Start Backend (watch) + Frontend (watch) together
	@echo "Starting Backend (watch) + Frontend (watch)..."
	@echo "Press Ctrl+C to stop all processes"
	@echo ""
	@make -j watch-backend watch-frontend

watch-backend: ## Start Backend with file watching (air)
	@cd backend && \
		BLOGGER_API_KEY=$(BLOGGER_API_KEY) \
		CONFIG_PATH=../$(CONFIG_PATH) \
		../tools/air

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

lint: ## Run ESLint + Go linter
	npm run lint
	@cd backend && ../../tools/golangci-lint run || echo "⚠ golangci-lint failed (optional)"

lint-fix: ## Fix linting issues (ESLint + gofmt)
	npm run lint:fix
	@cd backend && go fmt ./...

format: ## Format code (Prettier + gofmt)
	npm run format
	@cd backend && go fmt ./...

test: ## Run tests (Vitest + Go tests)
	npm run test
	cd backend && go test ./...

clean: ## Clean build artifacts
	rm -rf bin/ web/static/ coverage/ dist/
	@cd backend && rm -rf tmp/
	@echo "✓ Clean complete"

