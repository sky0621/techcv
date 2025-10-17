APP_NAME=api
APP_PATH=./cmd/api
COMPOSE ?= docker compose
SQLC ?= sqlc

GOLANGCI_LINT_VERSION ?= v2.5.0
TOOLS_BIN := $(CURDIR)/bin-tools
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
GOIMPORTS := $(TOOLS_BIN)/goimports

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-16s %s\n", $$1, $$2}'

.PHONY: run
run: ## Run the application in development mode
	go run $(APP_PATH)

.PHONY: build
build: ## Build the application binary
	go build -o bin/$(APP_NAME) $(APP_PATH)

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: tidy
tidy: ## Update go modules
	go mod tidy

.PHONY: download
download: ## Update go modules
	go mod download

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

.PHONY: lint-install
lint-install: ## Install golangci-lint locally
	@current="" ; \
	if [ -x "$(GOLANGCI_LINT)" ]; then \
		current=$$($(GOLANGCI_LINT) version 2>/dev/null | sed -E -n 's/.*version (v?[0-9.]+).*/\1/p'); \
		case $$current in \
			v*) ;; \
			"") ;; \
			*) current="v"$$current ;; \
		esac; \
	fi; \
	if [ "$$current" = "$(GOLANGCI_LINT_VERSION)" ]; then \
		echo "golangci-lint $(GOLANGCI_LINT_VERSION) already installed"; \
		else \
			echo "installing golangci-lint $(GOLANGCI_LINT_VERSION)"; \
			GOTOOLCHAIN=go1.25.0 GOCACHE=$(CURDIR)/.cache GOBIN=$(TOOLS_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
			rm -rf $(CURDIR)/.cache; \
		fi

.PHONY: fmt
fmt: ## Run gofmt on Go files
	go fmt ./...

.PHONY: imports-install
imports-install: ## Install goimports locally
	GOCACHE=$(CURDIR)/.cache GOBIN=$(TOOLS_BIN) go install golang.org/x/tools/cmd/goimports
	rm -rf $(CURDIR)/.cache

.PHONY: goimports
goimports: imports-install ## Run goimports on Go files
	$(GOIMPORTS) -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*' -not -path './.gomodcache/*')

.PHONY: imports
imports: goimports ## Backward compatibility alias

.PHONY: lint
lint: download lint-install ## Run golangci-lint
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
