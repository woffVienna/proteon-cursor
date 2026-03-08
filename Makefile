SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c

REPO_ROOT := $(CURDIR)

# -------- Included target groups --------
include $(REPO_ROOT)/tools/tooling.mk
include $(REPO_ROOT)/infra/k8s/make/infra.mk
include $(REPO_ROOT)/infra/k8s/make/prereqs.mk

# -------- Service discovery --------
# All service directories under services/*
SERVICES := $(shell find $(REPO_ROOT)/services -mindepth 1 -maxdepth 1 -type d -exec basename {} \; | sort)

# -------- Convenience --------
.PHONY: help
help:
	@echo "Targets:"
	@echo "  setup            - setup everything for local development"
	@echo "  create-service <name> - create a new service (e.g. make create-service api)"
	@echo "  generate         - run make generate in all services"
	@echo "  test             - run make test in all services"
	@echo "  build            - run make build in all services"
	@echo "  clean            - run make clean in all services"
	@$(MAKE) help-prereqs
	@$(MAKE) help-infra

# -------- Setup (one command) --------
.PHONY: setup
setup: check-prereqs tooling _work generate
	@echo "Setup complete."

# -------- go.work management --------
.PHONY: _work
_work:
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

# Public alias kept for compatibility with ENGINEERING.md, but hidden from help.
.PHONY: work
work: _work

# -------- Service creation --------
# Usage: make create-service <name>  (e.g. make create-service api)
SERVICE_NAME := $(firstword $(filter-out create-service,$(MAKECMDGOALS)))
.PHONY: create-service
create-service:
	@if [ -z "$(SERVICE_NAME)" ]; then \
		echo "Usage: make create-service <name>"; \
		exit 1; \
	fi
	@$(REPO_ROOT)/tools/scripts/create-service-structure.sh $(SERVICE_NAME)

# Consume the service name so make does not try to build it as a target
%:
	@:

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
