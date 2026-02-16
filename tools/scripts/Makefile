SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c

SERVICE ?= $(notdir $(CURDIR))
GO ?= go

REPO_ROOT := $(abspath $(CURDIR)/../..)
NODE_TOOLING_DIR := $(REPO_ROOT)/tools/node

REDOCLY := npx --prefix $(NODE_TOOLING_DIR) redocly
OAPI_CODEGEN ?= oapi-codegen

API_DIR := api
BUILD_DIR := .build
GEN_DIR := $(BUILD_DIR)/generated
BIN_DIR := $(BUILD_DIR)/bin

OPENAPI_SRC := $(API_DIR)/openapi.yml
OPENAPI_BUNDLE := $(GEN_DIR)/openapi.bundle.yml

HTTP_GEN_DIR := internal/adapters/http/generated/server
HTTP_GEN_FILE := $(HTTP_GEN_DIR)/openapi.gen.go

OAPI_SERVER_CFG := $(API_DIR)/oapi-codegen.server.yml

.PHONY: help
help:
	@echo "Service: $(SERVICE)"
	@echo "Targets: tooling tidy fmt lint test generate build run clean"

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: fmt
fmt:
	@find . -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -w

.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Skipping."; \
	fi

.PHONY: generate
generate: $(HTTP_GEN_FILE)

$(OPENAPI_BUNDLE): $(OPENAPI_SRC)
	@mkdir -p $(GEN_DIR)
	@if [[ ! -d "$(NODE_TOOLING_DIR)" ]]; then \
		echo "Error: tools/node not found at $(NODE_TOOLING_DIR)"; \
		exit 1; \
	fi
	@if [[ ! -d "$(NODE_TOOLING_DIR)/node_modules" ]]; then \
		echo "Error: node tooling not installed. Run: make tooling"; \
		exit 1; \
	fi
	@echo "Bundling OpenAPI -> $(OPENAPI_BUNDLE)"
	$(REDOCLY) bundle $(OPENAPI_SRC) --output $(OPENAPI_BUNDLE)

$(HTTP_GEN_FILE): $(OPENAPI_BUNDLE) $(OAPI_SERVER_CFG)
	@mkdir -p $(HTTP_GEN_DIR)
	@if command -v $(OAPI_CODEGEN) >/dev/null 2>&1; then \
		echo "Generating server code -> $(HTTP_GEN_FILE)"; \
		$(OAPI_CODEGEN) -config $(OAPI_SERVER_CFG) $(OPENAPI_BUNDLE); \
	else \
		echo "Error: oapi-codegen not installed."; \
		echo "Install: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"; \
		exit 1; \
	fi

.PHONY: dev
dev: generate
	$(GO) run ./cmd/$(SERVICE)

.PHONY: run
run:
	$(GO) run ./cmd/$(SERVICE)

.PHONY: test
test:
	$(GO) test ./...
	
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(SERVICE) ./cmd/$(SERVICE)

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
