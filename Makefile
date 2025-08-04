
# ERP Backend Development Makefile

# Configuration
PROTO_OUT_DIR = ${PWD}/gen
PROTO_DIT = ${PWD}/idl/proto/api_v1
DOCKER_COMPOSE = docker-compose
GO_CMD = go
APP_NAME = erp-backend
BUILD_DIR = build
COVERAGE_DIR = coverage

# Include protobuf rules
include Makefile.Protobuf.mk

# Default target
.DEFAULT_GOAL := help

## Development Commands

.PHONY: dev
dev: deps docker-dev-up ## Start development environment with hot-reload
	@echo "🚀 Starting ERP backend in development mode with hot-reload..."
	@echo "📊 Monitoring endpoints:"
	@echo "  - App: http://localhost:8080"
	@echo "  - NATS Monitoring: http://localhost:8222"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.toml; \
	else \
		echo "Installing air for hot-reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air -c .air.toml; \
	fi

.PHONY: dev-stop
dev-stop: ## Stop development environment
	@echo "🛑 Stopping development environment..."
	$(DOCKER_COMPOSE) down

.PHONY: docker-dev-up
docker-dev-up: ## Start only the infrastructure services (postgres, nats, redis)
	@echo "🐳 Starting development infrastructure..."
	$(DOCKER_COMPOSE) up -d postgres nats redis
	@echo "⏳ Waiting for services to be healthy..."
	$(DOCKER_COMPOSE) exec postgres pg_isready -U postgres -d erp_dev || sleep 5

.PHONY: docker-full-up
docker-full-up: ## Start full environment including app container
	@echo "🐳 Starting full development environment..."
	$(DOCKER_COMPOSE) --profile full up -d

## Build Commands

.PHONY: build
build: clean generate ## Build production binary
	@echo "🔨 Building production binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build \
		-ldflags="-w -s -X main.version=$$(git describe --tags --always --dirty)" \
		-o $(BUILD_DIR)/$(APP_NAME) \
		./cmd/all/main.go
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

.PHONY: build-local
build-local: generate ## Build local binary
	@echo "🔨 Building local binary..."
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build \
		-ldflags="-X main.version=$$(git describe --tags --always --dirty)" \
		-o $(BUILD_DIR)/$(APP_NAME) \
		./cmd/all/main.go
	@echo "✅ Local build complete: $(BUILD_DIR)/$(APP_NAME)"

## Testing Commands

.PHONY: test
test: docker-dev-up generate ## Run comprehensive test suite
	@echo "🧪 Running comprehensive test suite..."
	@mkdir -p $(COVERAGE_DIR)
	$(GO_CMD) test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	$(GO_CMD) tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "📊 Coverage report: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	$(GO_CMD) test -short -race ./...

.PHONY: test-integration
test-integration: docker-dev-up ## Run integration tests
	@echo "🧪 Running integration tests..."
	$(GO_CMD) test -v -race -tags=integration ./...

.PHONY: test-coverage
test-coverage: ## Generate test coverage report
	@echo "📊 Generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	$(GO_CMD) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report: $(COVERAGE_DIR)/coverage.html"

## Code Generation Commands

.PHONY: generate
generate: models annotations mockery ## Run all code generation
	@echo "✅ Code generation complete"

.PHONY: models
models: ## Generate GORM models from database schema
	@echo "🔄 Generating GORM models..."
	@cd cmd/gen-models && $(GO_CMD) run main.go 

.PHONY: annotations
annotations: ## Generate annotations
	@echo "🔄 Generating annotations..."
	@cd cmd/gen-annotations && $(GO_CMD) run main.go

.PHONY: mockery
mockery: ## Generate mocks for interfaces
	@echo "🔄 Generating mocks..."
	@if command -v mockery >/dev/null 2>&1; then \
		mockery; \
	else \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
		mockery; \
	fi

.PHONY: openapi
openapi: ## Generate OpenAPI specifications
	@echo "🔄 Generating OpenAPI specs..."
	@mkdir -p gen/openapi
	@cd cmd/all && $(GO_CMD) run main.go --generate-openapi
	@echo "✅ OpenAPI specs generated in gen/openapi/"

## Protocol Buffers

.PHONY: protoc
protoc: ## Generate protobuf code
	@echo "🔄 Generating protobuf code..."
	protoc -Iidl/proto/api_v1 --go_out=${PROTO_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${PROTO_OUT_DIR} \
	--go-grpc_opt=paths=source_relative idl/proto/api_v1/*.proto

## Code Quality Commands

.PHONY: lint
lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "🎨 Formatting code..."
	$(GO_CMD) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "🔍 Running go vet..."
	$(GO_CMD) vet ./...

.PHONY: check
check: fmt vet lint test-unit ## Run all code quality checks

## Database Commands

.PHONY: db-reset
db-reset: ## Reset development database
	@echo "🗄️ Resetting development database..."
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS erp_dev;"
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -c "CREATE DATABASE erp_dev;"
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d erp_dev -f /docker-entrypoint-initdb.d/init.sql

.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "🗄️ Running database migrations..."
	@cd cmd/all && $(GO_CMD) run main.go --migrate

## Utility Commands

.PHONY: deps
deps: ## Install/update dependencies
	@echo "📦 Installing dependencies..."
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy

.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	$(GO_CMD) clean -cache -testcache -modcache

.PHONY: docker
docker: ## Build and push Docker image
	@echo "🐳 Building Docker image..."
	docker build -t jmiranda0521/erp:$(if $(tag),$(tag),latest) .
	@if [ "$(push)" = "true" ]; then \
		echo "📤 Pushing Docker image..."; \
		docker push jmiranda0521/erp:$(if $(tag),$(tag),latest); \
	fi

## Legacy Commands (maintained for compatibility)

.PHONY: accounting
accounting: ## Legacy: Run accounting service
	@cd accounting/cmd && $(GO_CMD) run main.go

.PHONY: app
app: ## Legacy: Run app service  
	@cd cmd/app && $(GO_CMD) run main.go

.PHONY: start
start: ## Legacy: Start application
	@cd cmd/all && $(GO_CMD) run main.go

.PHONY: consul
consul: ## Legacy: Start Consul agent
	consul agent -dev -ui -node member

.PHONY: test-short
test-short: ## Legacy: Run short tests
	$(GO_CMD) test -short ./...

.PHONY: doc
doc: ## Legacy: Generate documentation
	docker run --rm -v ${PWD}/documents:/documents asciidoctor/docker-asciidoctor sample.adoc

## Database Backup (Legacy)
PG_URL = postgresql://postgres:12ab34cd56ef@10.0.0.151:5432/erp_dev

.PHONY: backup-test-db
backup-test-db: ## Legacy: Backup test database
	pg_dump ${PG_URL} --format plain --data-only --verbose --file "db/data.sql" --table entities \
	--table actions --table currencies --table party_types --table states --table unit_of_measures \
	--table unit_of_measure_translations 
	pg_dump ${PG_URL} --format plain --schema-only --verbose --file "db/schema.sql"
	cd db && copy schema.sql + data.sql + custom-data.sql init.sql

## Help

.PHONY: help
help: ## Show this help message
	@echo "📋 ERP Backend Development Commands"
	@echo "=================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make dev          # Start development environment with hot-reload" 
	@echo "  make test         # Run comprehensive tests"
	@echo "  make build        # Build production binary"
