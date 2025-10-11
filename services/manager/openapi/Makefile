NPM_BIN := npm
REDOCLY := node_modules/.bin/redocly
SPEC_INDEX := spec/index.yaml
OUTPUT := openapi.yaml

.PHONY: install-redocly
install-redocly: $(REDOCLY)

$(REDOCLY): package.json
	$(NPM_BIN) install

.PHONY: bundle-openapi
bundle-openapi: $(REDOCLY)
	$(REDOCLY) bundle $(SPEC_INDEX) --output $(OUTPUT)

.PHONY: clean
clean:
	rm -rf node_modules
