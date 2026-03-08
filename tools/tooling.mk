# -------- Tooling locations --------
NODE_TOOLING_DIR := $(REPO_ROOT)/tools/node
TOOLS_BIN := $(REPO_ROOT)/tools/bin

# Ensure locally installed Go tools are picked up
export PATH := $(TOOLS_BIN):$(PATH)

# -------- Tooling --------
.PHONY: tooling
tooling: tooling-node tooling-go

.PHONY: tooling-node
tooling-node:
	@if [[ -f "$(NODE_TOOLING_DIR)/package-lock.json" ]]; then \
		echo "Installing node tooling via npm ci in $(NODE_TOOLING_DIR)"; \
		( cd "$(NODE_TOOLING_DIR)" && npm ci ); \
	elif [[ -f "$(NODE_TOOLING_DIR)/package.json" ]]; then \
		echo "Installing node tooling via npm install in $(NODE_TOOLING_DIR)"; \
		( cd "$(NODE_TOOLING_DIR)" && npm install ); \
	else \
		echo "Error: missing tools/node/package.json at $(NODE_TOOLING_DIR)"; \
		exit 1; \
	fi

# Installs Go tools locally (no global installs)
.PHONY: tooling-go
tooling-go:
	@mkdir -p "$(TOOLS_BIN)"
	@echo "Installing Go tools into $(TOOLS_BIN)"
	@GOBIN="$(TOOLS_BIN)" go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@GOBIN="$(TOOLS_BIN)" go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
