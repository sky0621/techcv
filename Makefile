SERVICE_DIRS := $(sort $(wildcard services/*/))
LAYER_TOKENS := be backend fe frontend
LAYER_ALIAS_be := backend
LAYER_ALIAS_backend := backend
LAYER_ALIAS_fe := frontend
LAYER_ALIAS_frontend := frontend

OPENAPI_TARGETS := install-redocly bundle-openapi clean
.PHONY: $(OPENAPI_TARGETS:%=openapi-%)

define DISPATCH_SERVICE_TARGET
$(eval __goal := $(1))
$(eval __alias := $(firstword $(subst -, ,$(__goal))))
$(eval __rest := $(patsubst $(__alias)-%,%,$(__goal)))
$(if $(strip $(__rest)),,$(error Invalid target '$(__goal)' (missing layer segment)))
$(eval __layer_alias := $(firstword $(subst -, ,$(__rest))))
$(if $(strip $(__layer_alias)),,$(error Invalid target '$(__goal)' (missing layer alias)))
$(if $(findstring -,$(__rest)),,$(error Invalid target '$(__goal)' (missing command segment)))
$(eval __target := $(patsubst $(__layer_alias)-%,%,$(__rest)))
$(eval __matches := $(filter services/$(__alias)%/,$(SERVICE_DIRS)))
$(if $(strip $(__matches)),,$(error Unknown service alias '$(__alias)'))
$(if $(filter-out 1,$(words $(__matches))),$(error Ambiguous service alias '$(__alias)' matches $(__matches)))
$(eval __match := $(firstword $(__matches)))
$(eval __service := $(patsubst services/%/,%,$(__match)))
$(eval __layer := $(LAYER_ALIAS_$(__layer_alias)))
$(if $(strip $(__layer)),,$(error Unknown layer alias '$(__layer_alias)'))
$(MAKE) -C services/$(__service)/$(__layer) $(__target)
endef

define DEFINE_SERVICE_GOAL
$(1):
	$$(call DISPATCH_SERVICE_TARGET,$(1))

.PHONY: $(1)
endef

IS_SERVICE_GOAL = $(strip $(foreach layer,$(LAYER_TOKENS),$(findstring -$(layer)-,$(1))))

# manager-specific openapi commands remain available
openapi-%:
	$(MAKE) -C services/manager/openapi $(@:openapi-%=%)

ifneq ($(MAKECMDGOALS),)
$(foreach goal,$(MAKECMDGOALS), \
  $(if $(call IS_SERVICE_GOAL,$(goal)), \
    $(eval $(call DEFINE_SERVICE_GOAL,$(goal))) \
  ) \
)
endif
