APP_NAME=api
APP_PATH=./cmd/api
COMPOSE ?= docker compose
SQLC ?= sqlc

GOLANGCI_LINT ?= golangci-lint
GOIMPORTS ?= goimports

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-16s %s\n", $$1, $$2}'

# Local development
.PHONY: run
run: ## Run the application in development mode
	go run $(APP_PATH)

.PHONY: build
build: ## Build the application binary
	go build -o bin/$(APP_NAME) $(APP_PATH)

.PHONY: test
test: ## Run unit tests
	go test ./...

# Dependency management
.PHONY: tidy
tidy: ## Update go modules
	go mod tidy

.PHONY: download
download: ## Update go modules
	go mod download

# Database and code generation
.PHONY: migrate
migrate: ## Apply schema changes to the local database using sqldef
	$(COMPOSE) run --rm sqldef

.PHONY: gen-sqlc
gen-sqlc: ## Generate database access layer code with sqlc
	$(SQLC) generate

.PHONY: gen-api
gen-api: ## Generate Go sources from OpenAPI specification
	GOCACHE=$(CURDIR)/.cache go run ./cmd/openapi -spec ../openapi/openapi.yaml -out ./internal/interface/http/openapi/openapi.gen.go
	rm -rf $(CURDIR)/.cache

# Formatting
.PHONY: fmt
fmt: ## Run gofmt on Go files
	go fmt ./...

.PHONY: goimports
goimports: ## Run goimports on Go files
	$(GOIMPORTS) -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*' -not -path './.gomodcache/*')

.PHONY: imports
imports: goimports ## Backward compatibility alias

# Quality gates
.PHONY: lint
lint: download ## Run golangci-lint
	GOCACHE=$(CURDIR)/.cache GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/golangci-lint $(GOLANGCI_LINT) run
	rm -rf $(CURDIR)/.cache

.PHONY: check
check:
	make tidy
	make download
	make fmt
	make imports
	make lint
	make build
	make test
