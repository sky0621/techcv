APP_NAME=api
APP_PATH=./cmd/api
GOLANGCI_LINT_VERSION ?= v1.58.2
TOOLS_BIN := $(PWD)/bin-tools
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint

.PHONY: help run build test tidy lint generate-openapi golangci-lint-install fmt

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

golangci-lint-install: ## Install golangci-lint locally
	GOBIN=$(TOOLS_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

fmt: ## Run gofmt on Go files
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*')

lint: golangci-lint-install ## Run golangci-lint
	$(GOLANGCI_LINT) run

generate-openapi: ## Generate Go sources from OpenAPI specification
	GOCACHE=$(PWD)/.cache go run ./cmd/openapi -spec ../openapi/openapi.yaml -out ./internal/interface/http/openapi/openapi.gen.go
	rm -rf $(PWD)/.cache
