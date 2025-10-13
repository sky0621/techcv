BACKEND_DIR := services/manager/backend
OPENAPI_DIR := services/manager/openapi
FRONTEND_DIR := services/manager/frontend

BACKEND_TARGETS := run build test tidy lint generate-openapi help
OPENAPI_TARGETS := install-redocly bundle-openapi clean
FRONTEND_TARGETS := install vite-install dev build preview lint test

.PHONY: $(addprefix backend-,$(BACKEND_TARGETS)) $(addprefix openapi-,$(OPENAPI_TARGETS)) $(addprefix frontend-,$(FRONTEND_TARGETS))

$(addprefix backend-,$(BACKEND_TARGETS)):
	$(MAKE) -C $(BACKEND_DIR) $(@:backend-%=%)

$(addprefix openapi-,$(OPENAPI_TARGETS)):
	$(MAKE) -C $(OPENAPI_DIR) $(@:openapi-%=%)

$(addprefix frontend-,$(FRONTEND_TARGETS)):
	$(MAKE) -C $(FRONTEND_DIR) $(@:frontend-%=%)
