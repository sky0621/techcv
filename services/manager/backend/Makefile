APP_NAME=api
APP_PATH=./cmd/api
COMPOSE ?= docker compose
COMPOSE_FILE ?= $(CURDIR)/docker-compose.yml
SQLC ?= sqlc

SQLDEF_BIN ?= $(CURDIR)/bin-tools/mysqldef
SQLDEF_PKG := github.com/sqldef/sqldef/cmd/mysqldef

DB_HOST ?= 127.0.0.1
DB_PORT ?= 3306
DB_NAME ?= manager
DB_USER ?= manager
DB_PASSWORD ?= manager
DB_ROOT_PASSWORD ?= root
SCHEMA_FILE ?= db/schema.sql

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
.PHONY: db-up
db-up: ## Start the local MySQL container
	DB_NAME=$(DB_NAME) \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_ROOT_PASSWORD=$(DB_ROOT_PASSWORD) \
	DB_PORT=$(DB_PORT) \
	$(COMPOSE) -f $(COMPOSE_FILE) up -d db

.PHONY: db-down
db-down: ## Stop the local MySQL container
	DB_NAME=$(DB_NAME) \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_ROOT_PASSWORD=$(DB_ROOT_PASSWORD) \
	DB_PORT=$(DB_PORT) \
	$(COMPOSE) -f $(COMPOSE_FILE) stop db

.PHONY: migrate
migrate: $(SQLDEF_BIN) ## Apply schema changes to the local database using mysqldef
	MYSQL_PWD=$(DB_PASSWORD) $(SQLDEF_BIN) \
		-h $(DB_HOST) \
		-P $(DB_PORT) \
		-u $(DB_USER) \
		--file $(SCHEMA_FILE) \
		$(DB_NAME)

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

.PHONY: imports
imports: ## Run goimports on Go files
	go tool golang.org/x/tools/cmd/goimports -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*' -not -path './.gomodcache/*')

# Quality gates
.PHONY: lint
lint: ## Run golangci-lint
	GOCACHE=$(CURDIR)/.cache GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/golangci-lint $(CURDIR)/bin-tools/golangci-lint run
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

$(SQLDEF_BIN):
	@mkdir -p $(dir $(SQLDEF_BIN))
	@version=$$(cd $(CURDIR) && go list -f '{{.Version}}' -m $(SQLDEF_PKG) 2>/dev/null || echo latest); \
	echo "Installing mysqldef ($$version)"; \
	GO111MODULE=on GOBIN=$(CURDIR)/bin-tools go install $(SQLDEF_PKG)@$$version
