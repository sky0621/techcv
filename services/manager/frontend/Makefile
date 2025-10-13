NPM_BIN ?= npm

.PHONY: install vite-install dev build preview lint test

install: ## Install node modules
	$(NPM_BIN) install

vite-install: install ## Ensure vite is installed (via npm dependencies)
	@echo "Vite available via node_modules/.bin/vite"

dev: ## Start Vite dev server
	npx vite

build: ## Build production bundle
	npx tsc -b
	npx vite build

preview: ## Preview production bundle
	npx vite preview

lint: ## Run eslint checks
	npx eslint src --ext .ts,.tsx

test: ## Run Jest tests
	npx jest
