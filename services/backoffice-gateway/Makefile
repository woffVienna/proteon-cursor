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

# Gateway uses reverse proxy for downstream routes; no server codegen needed.
# Generate only bundles the OpenAPI spec for Swagger UI.
.PHONY: generate
generate: $(OPENAPI_BUNDLE)

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

# -------- Containerisation --------

IMAGE_REPO ?= proteon
IMAGE_TAG ?= dev
DOCKERFILE ?= Dockerfile
# Keep runtime/build paths based on SERVICE (= folder name), but allow
# docker image naming to follow deployed service naming.
IMAGE_NAME ?= $(SERVICE)-service
IMAGE := $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: containerise
containerise: generate
	@echo "Building docker image -> $(IMAGE)"
	@docker build \
		-f "$(REPO_ROOT)/services/$(SERVICE)/$(DOCKERFILE)" \
		-t "$(IMAGE)" \
		"$(REPO_ROOT)"

