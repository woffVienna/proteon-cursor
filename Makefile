SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c

REPO_ROOT := $(CURDIR)

# -------- Tooling locations --------
NODE_TOOLING_DIR := $(REPO_ROOT)/tools/node
TOOLS_BIN := $(REPO_ROOT)/tools/bin

# Ensure locally installed Go tools are picked up
export PATH := $(TOOLS_BIN):$(PATH)

# -------- Service discovery --------
# All service directories under services/*
SERVICES := $(shell find $(REPO_ROOT)/services -mindepth 1 -maxdepth 1 -type d -exec basename {} \; | sort)

# -------- Convenience --------
.PHONY: help
help:
	@echo "Targets:"
	@echo "  setup            - setup everything for local development"
	@echo "  tooling          - install all tooling (node + go tools)"
	@echo "  tooling-node     - install node tooling in tools/node"
	@echo "  tooling-go       - install go tools into tools/bin"
	@echo "  work             - init/update go.work to include all services + libs/platform"
	@echo "  generate         - run make generate in all services"
	@echo "  test             - run make test in all services"
	@echo "  build            - run make build in all services"
	@echo "  clean            - run make clean in all services"

# -------- Setup (one command) --------
.PHONY: setup
setup: tooling work generate
	@echo "Setup complete."

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

# -------- go.work management --------
.PHONY: work
work:
	@if [[ ! -f "$(REPO_ROOT)/go.work" ]]; then \
		echo "Initializing go.work"; \
		go work init; \
	fi
	@if [[ -d "$(REPO_ROOT)/libs/platform" ]]; then \
		echo "Adding libs/platform to go.work"; \
		go work use ./libs/platform; \
	fi
	@for s in $(SERVICES); do \
		echo "Adding services/$$s to go.work"; \
		go work use ./services/$$s; \
	done
	@echo "go.work updated."

# -------- Fan-out targets to services --------
.PHONY: generate
generate:
	@for s in $(SERVICES); do \
		echo "==> generating $$s"; \
		$(MAKE) -C services/$$s generate; \
	done

.PHONY: test
test:
	@for s in $(SERVICES); do \
		echo "==> testing $$s"; \
		$(MAKE) -C services/$$s test; \
	done

.PHONY: build
build:
	@for s in $(SERVICES); do \
		echo "==> building $$s"; \
		$(MAKE) -C services/$$s build; \
	done

.PHONY: clean
clean:
	@for s in $(SERVICES); do \
		echo "==> cleaning $$s"; \
		$(MAKE) -C services/$$s clean; \
	done

# -------- CI/CD checks --------
.PHONY: verify-generated
verify-generated: generate
	@echo "Verifying generated code is up to date..."
	@if ! git diff --quiet -- services/*/internal/**/generated/**; then \
		echo ""; \
		echo "ERROR: Generated files are out of date."; \
		echo "Run 'make generate' and commit the changes."; \
		echo ""; \
		git --no-pager diff -- services/*/internal/**/generated/**; \
		exit 1; \
	fi
	@echo "OK: generated code is up to date."

.PHONY: check
check: verify-generated test
