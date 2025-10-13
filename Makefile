SERVICES := manager publisher administrator
BACKEND_TARGETS := run build test tidy fmt lint gen-api help
FRONTEND_TARGETS := install vite-install dev build preview lint test
OPENAPI_TARGETS := install-redocly bundle-openapi clean

.PHONY: $(foreach svc,$(SERVICES),$(foreach tgt,$(BACKEND_TARGETS),$(svc)-be-$(tgt))) \
        $(foreach svc,$(SERVICES),$(foreach tgt,$(FRONTEND_TARGETS),$(svc)-fe-$(tgt))) \
        openapi-$(OPENAPI_TARGETS)

$(foreach svc,$(SERVICES),$(foreach tgt,$(BACKEND_TARGETS),$(svc)-be-$(tgt))):
	$(eval SERVICE := $(word 1,$(subst -, ,$@)))
	$(eval TARGET := $(word 3,$(subst -, ,$@)))
	$(MAKE) -C services/$(SERVICE)/backend $(TARGET)

$(foreach svc,$(SERVICES),$(foreach tgt,$(FRONTEND_TARGETS),$(svc)-fe-$(tgt))):
	$(eval SERVICE := $(word 1,$(subst -, ,$@)))
	$(eval TARGET := $(word 3,$(subst -, ,$@)))
	$(MAKE) -C services/$(SERVICE)/frontend $(TARGET)

# manager-specific openapi commands remain available
openapi-%:
	$(MAKE) -C services/manager/openapi $(@:openapi-%=%)
