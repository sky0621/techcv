.PHONY: run test fmt vet tidy healthcheck

APP_ENV ?= local
APP_PORT ?= 8080
SQLITE_PATH ?= ./tmp/manager.db
GOCACHE ?= $(CURDIR)/.cache/go-build
HEALTHCHECK_URL ?= http://127.0.0.1:$(APP_PORT)/health

run:
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

healthcheck:
	@body=$$(curl -fsS "$(HEALTHCHECK_URL)"); \
	echo "$$body"; \
	echo "$$body" | grep -q '"status":"ok"' || (echo "healthcheck failed: status is not ok" >&2; exit 1)
