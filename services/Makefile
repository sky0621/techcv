APP_NAME=api
APP_PATH=./cmd/api

.PHONY: help run build test tidy lint

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
