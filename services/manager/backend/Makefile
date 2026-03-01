.PHONY: run test fmt vet tidy lint check healthcheck sqlc-generate openapi-generate db-migrate

APP_ENV ?= local
APP_PORT ?= 8080
SQLITE_PATH ?= ./tmp/manager.db
GOCACHE ?= $(CURDIR)/.cache/go-build
GOLANGCI_LINT_CACHE ?= $(CURDIR)/.cache/golangci-lint
HEALTHCHECK_URL ?= http://127.0.0.1:$(APP_PORT)/health

run: db-migrate
	mkdir -p "$(GOCACHE)"
	APP_ENV=$(APP_ENV) APP_PORT=$(APP_PORT) SQLITE_PATH=$(SQLITE_PATH) GOCACHE=$(GOCACHE) go run ./cmd/api

test:
	mkdir -p "$(GOCACHE)"
	GOCACHE=$(GOCACHE) go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -type f)

vet:
	mkdir -p "$(GOCACHE)"
	GOCACHE=$(GOCACHE) go vet ./...

tidy:
	mkdir -p "$(GOCACHE)"
	GOCACHE=$(GOCACHE) go mod tidy

lint:
	@command -v golangci-lint >/dev/null 2>&1 || (echo "golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8" >&2; exit 1)
	mkdir -p "$(GOCACHE)" "$(GOLANGCI_LINT_CACHE)"
	GOCACHE=$(GOCACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) golangci-lint run ./...

check:
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) tidy
	$(MAKE) lint
	$(MAKE) test

sqlc-generate:
	@command -v sqlc >/dev/null 2>&1 || (echo "sqlc not found. Install: go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.24.0" >&2; exit 1)
	sqlc generate

openapi-generate:
	mkdir -p "$(GOCACHE)"
	mkdir -p internal/interface/http/openapi/gen
	GOCACHE=$(GOCACHE) go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen.yaml openapi/manager-api.yaml

db-migrate:
	@if [ "$(SQLITE_PATH)" = ":memory:" ]; then \
		echo "skip db-migrate for SQLITE_PATH=:memory:"; \
		exit 0; \
	fi
	@command -v sqlite3def >/dev/null 2>&1 || (echo "sqlite3def not found. Install: go install github.com/sqldef/sqldef/cmd/sqlite3def@latest" >&2; exit 1)
	mkdir -p "$(dir $(SQLITE_PATH))"
	touch "$(SQLITE_PATH)"
	sqlite3def "$(SQLITE_PATH)" --file sql/schema.sql

healthcheck:
	@body=$$(curl -fsS "$(HEALTHCHECK_URL)"); \
	echo "$$body"; \
	echo "$$body" | grep -q '"status":"ok"' || (echo "healthcheck failed: status is not ok" >&2; exit 1)
