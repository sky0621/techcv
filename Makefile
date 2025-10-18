# 動的ディスパッチのためにサービスディレクトリ (services/<name>/) を列挙する。
SERVICE_DIRS := $(sort $(wildcard services/*/))
# ディスパッチャが認識する正規レイヤートークン。
LAYER_TOKENS := be backend fe frontend
# 短いエイリアスを正規のバックエンドレイヤー名に対応付ける。
LAYER_ALIAS_be := backend
LAYER_ALIAS_backend := backend
# 短いエイリアスを正規のフロントエンドレイヤー名に対応付ける。
LAYER_ALIAS_fe := frontend
LAYER_ALIAS_frontend := frontend

# 標準的な OpenAPI 管理ターゲットを公開する。
OPENAPI_TARGETS := install-redocly bundle-openapi clean
.PHONY: $(OPENAPI_TARGETS:%=openapi-%)

# ディスパッチヘルパーはサービスレイヤーコマンド (alias-layer-target) を展開して make を呼び出す。
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

# 呼び出し側が要求したサービスレイヤーのゴールに対する PHONY ルールを生成する。
define DEFINE_SERVICE_GOAL
$(1):
	$$(call DISPATCH_SERVICE_TARGET,$(1))

.PHONY: $(1)
endef

# ゴールが <service>-<layer>-<target> 形式かどうかを判定する述語。
IS_SERVICE_GOAL = $(strip $(foreach layer,$(LAYER_TOKENS),$(findstring -$(layer)-,$(1))))

# manager 専用の openapi コマンドを引き続き利用可能にする。
openapi-%:
	$(MAKE) -C services/manager/openapi $(@:openapi-%=%)

# 要求された各ゴールに対してサービスレイヤーのルールを動的に定義する。
ifneq ($(MAKECMDGOALS),)
$(foreach goal,$(MAKECMDGOALS), \
  $(if $(call IS_SERVICE_GOAL,$(goal)), \
    $(eval $(call DEFINE_SERVICE_GOAL,$(goal))) \
  ) \
)
endif
