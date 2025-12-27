APP_NAME=api
APP_PATH=./cmd/api
COMPOSE ?= docker compose
SQLC ?= sqlc

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
	go tool sqlc generate

.PHONY: gen-api
gen-api: ## Generate Go sources from OpenAPI specification
	go tool oapi-codegen --config ./oapi-codegen.yaml ../openapi/openapi.yaml

# Formatting
.PHONY: fmt
fmt: ## Run gofmt on Go files
	go fmt ./...

.PHONY: imports
imports: ## Run goimports on Go files
	go tool goimports -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*' -not -path './.gomodcache/*')

# Quality gates
.PHONY: lint
lint: ## Run golangci-lint
	go tool golangci-lint run

.PHONY: check
check:
	make tidy
	make download
	make fmt
	make imports
	make lint
	make build
	make test
