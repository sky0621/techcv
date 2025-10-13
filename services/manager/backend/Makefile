APP_NAME=api
APP_PATH=./cmd/api
GOLANGCI_LINT_VERSION ?= v1.60.2
TOOLS_BIN := $(PWD)/bin-tools
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
GOIMPORTS := $(TOOLS_BIN)/goimports
GOIMPORTS_VERSION ?= latest

.PHONY: help run build test tidy download fmt lint gen-api lint-install goimports goimports-install imports imports-install

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-16s %s\n", $$1, $$2}'

run: ## Run the application in development mode
	go run $(APP_PATH)

build: ## Build the application binary
	go build -o bin/$(APP_NAME) $(APP_PATH)

test: ## Run unit tests
	go test ./...

tidy: ## Update go modules
	go mod tidy

download: ## Update go modules
	go mod download

lint-install: ## Install golangci-lint locally
	GOBIN=$(TOOLS_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

fmt: ## Run gofmt on Go files
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*')

goimports-install: ## Install goimports locally
	GOBIN=$(TOOLS_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

goimports: goimports-install ## Run goimports on Go files
	$(GOIMPORTS) -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*')

imports: goimports ## Backward compatibility alias

imports-install: goimports-install ## Backward compatibility alias

lint: download lint-install ## Run golangci-lint
	GOCACHE=$(PWD)/.cache GOLANGCI_LINT_CACHE=$(PWD)/.cache/golangci-lint $(GOLANGCI_LINT) run
	rm -rf $(PWD)/.cache

gen-api: ## Generate Go sources from OpenAPI specification
	GOCACHE=$(PWD)/.cache go run ./cmd/openapi -spec ../openapi/openapi.yaml -out ./internal/interface/http/openapi/openapi.gen.go
	rm -rf $(PWD)/.cache
