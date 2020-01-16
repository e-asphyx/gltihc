DIST_DIR = ./dist
SRC_DIR = ./pkg
STATIC_DIR = ./static

GO_SRC := $(SRC_DIR)/main.go
STATIC_SRC := \
	./static/index.html \
	./static/style.css \
	./node_modules/normalize.css/normalize.css

BIN = $(DIST_DIR)/gltihc.wasm
STATIC = $(addprefix $(DIST_DIR)/, $(notdir $(STATIC_SRC)))

.PHONY: all typescript

all: $(BIN) $(STATIC) typescript

$(BIN): $(GO_SRC)
	GOOS=js GOARCH=wasm go build -o $@ $(SRC_DIR)

define STATIC_tpl
$$(DIST_DIR)/$(notdir $(1)): $(1)
	mkdir -p $$(dir $$@)
	cp -a $$< $$@
endef

$(foreach s, $(STATIC_SRC), $(eval $(call STATIC_tpl, $(s))))

typescript:
	npx tsc -p .