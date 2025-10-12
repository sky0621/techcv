BACKEND_DIR := services/manager/backend
OPENAPI_DIR := services/manager/openapi

BACKEND_TARGETS := run build test tidy lint generate-openapi help
OPENAPI_TARGETS := install-redocly bundle-openapi clean

.PHONY: $(addprefix backend-,$(BACKEND_TARGETS)) $(addprefix openapi-,$(OPENAPI_TARGETS))

$(addprefix backend-,$(BACKEND_TARGETS)):
	$(MAKE) -C $(BACKEND_DIR) $(@:backend-%=%)

$(addprefix openapi-,$(OPENAPI_TARGETS)):
	$(MAKE) -C $(OPENAPI_DIR) $(@:openapi-%=%)
