APP_NAME=api
APP_PATH=./cmd/api
COMPOSE ?= docker compose
SQLC ?= sqlc

.PHONY: help run build test tidy lint migrate generate-sqlc generate-openapi

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

lint: ## Run gofmt to ensure formatting
	gofmt -l $(shell find . -name '*.go' -not -path './vendor/*')

migrate: ## Apply schema changes to the local database using sqldef
	$(COMPOSE) run --rm sqldef

generate-sqlc: ## Generate database access layer code with sqlc
	$(SQLC) generate

generate-openapi: ## Generate Go sources from OpenAPI specification
	GOCACHE=$(PWD)/.cache go run ./cmd/openapi -spec ../openapi/openapi.yaml -out ./internal/interface/http/openapi/openapi.gen.go
	rm -rf $(PWD)/.cache
