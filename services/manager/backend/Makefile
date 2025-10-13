APP_NAME=api
APP_PATH=./cmd/api
GOLANGCI_LINT_VERSION ?= v2.5.0
TOOLS_BIN := $(CURDIR)/bin-tools
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

fmt: ## Run gofmt on Go files
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*')

goimports-install: ## Install goimports locally
	GOBIN=$(TOOLS_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

goimports: goimports-install ## Run goimports on Go files
	$(GOIMPORTS) -w $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin-tools/*')

imports: goimports ## Backward compatibility alias

imports-install: goimports-install ## Backward compatibility alias

lint: download lint-install ## Run golangci-lint
	GOCACHE=$(CURDIR)/.cache GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/golangci-lint $(GOLANGCI_LINT) run
	rm -rf $(CURDIR)/.cache

gen-api: ## Generate Go sources from OpenAPI specification
	GOCACHE=$(CURDIR)/.cache go run ./cmd/openapi -spec ../openapi/openapi.yaml -out ./internal/interface/http/openapi/openapi.gen.go
	rm -rf $(CURDIR)/.cache
