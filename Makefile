SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c

REPO_ROOT := $(CURDIR)
HOST_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOST_ARCH_RAW := $(shell uname -m)
ifeq ($(HOST_ARCH_RAW),x86_64)
HOST_ARCH := amd64
else ifeq ($(HOST_ARCH_RAW),aarch64)
HOST_ARCH := arm64
else
HOST_ARCH := $(HOST_ARCH_RAW)
endif
HOST_PLATFORM := $(HOST_OS)-$(HOST_ARCH)

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
	@echo "  check-prereqs    - verify required local stack binaries are installed"
	@echo "  setup-darwin-arm64-prereqs - install local stack prerequisites on darwin-arm64"
	@echo "  tooling          - install all tooling (node + go tools)"
	@echo "  tooling-node     - install node tooling in tools/node"
	@echo "  tooling-go       - install go tools into tools/bin"
	@echo "  work             - init/update go.work to include all services + libs/platform"
	@echo "  create-service <name> - create a new service (e.g. make create-service api)"
	@echo "  generate         - run make generate in all services"
	@echo "  test             - run make test in all services"
	@echo "  build            - run make build in all services"
	@echo "  clean            - run make clean in all services"
	@echo "  cluster-up       - create local k3d cluster"
	@echo "  cluster-down     - delete local k3d cluster"
	@echo "  ns-up            - create kubernetes namespace"
	@echo "  ns-down          - delete kubernetes namespace"
	@echo "  deps-install     - install postgres + nats via helm"
	@echo "  deps-uninstall   - uninstall postgres + nats helm releases"
	@echo "  image-load       - build image and import into k3d (SERVICE=...)"
	@echo "  deploy           - deploy service helm chart (SERVICE=...)"
	@echo "  deploy-all       - deploy all local helm-charted services"
	@echo "  wait-deps        - wait for postgres + nats readiness"
	@echo "  wait-services    - wait for service rollout readiness"
	@echo "  wait-ingress     - wait for local ingress HTTP readiness"
	@echo "  stack-up         - full local stack bring-up (k3d + helm)"
	@echo "  stack-down       - uninstall local helm releases"

# -------- Setup (one command) --------
.PHONY: setup
setup: check-prereqs tooling work generate
	@echo "Setup complete."

.PHONY: check-prereqs
check-prereqs:
	@missing=""; \
	for cmd in docker k3d kubectl helm curl; do \
		if ! command -v $$cmd >/dev/null 2>&1; then \
			missing="$$missing $$cmd"; \
		fi; \
	done; \
	if [[ -n "$$missing" ]]; then \
		echo "Missing required tools:$$missing"; \
		echo "Run: make setup-$(HOST_PLATFORM)-prereqs"; \
		exit 1; \
	fi
	@echo "All required local stack prerequisites are installed."

.PHONY: setup-darwin-arm64-prereqs
setup-darwin-arm64-prereqs:
	@if ! command -v brew >/dev/null 2>&1; then \
		echo "Homebrew is required. Install it from https://brew.sh and re-run this target."; \
		exit 1; \
	fi
	@brew install k3d kubectl helm
	@brew install --cask docker
	@echo "Installed prerequisites for darwin-arm64."
	@echo "Start Docker Desktop once (for the Docker daemon): open -a Docker"

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

# -------- Local kubernetes orchestration (k3d + helm) --------
K3D_CLUSTER_NAME := proteon
K8S_NAMESPACE := proteon-dev
K3D_CLUSTER_SCRIPT := $(REPO_ROOT)/infra/k8s/local/k3d/cluster.sh

HELM_VALUES_DIR := $(REPO_ROOT)/infra/k8s/local/helm
HELM_CHARTS_DIR := $(REPO_ROOT)/infra/k8s/charts
POSTGRES_VALUES_FILE := $(HELM_VALUES_DIR)/postgres-values.yaml
NATS_VALUES_FILE := $(HELM_VALUES_DIR)/nats-values.yaml

POSTGRES_RELEASE := postgresql
NATS_RELEASE := nats
DEPLOYABLE_SERVICES := identity

IMAGE_REPO ?= proteon
IMAGE_TAG ?= dev
WAIT_TIMEOUT ?= 180s
WAIT_RETRY_SECONDS ?= 2
HEALTHCHECK_URL ?= http://localhost:8080/v1/health

.PHONY: cluster-up
cluster-up:
	@bash "$(K3D_CLUSTER_SCRIPT)" up

.PHONY: cluster-down
cluster-down:
	@bash "$(K3D_CLUSTER_SCRIPT)" down

.PHONY: ns-up
ns-up:
	@kubectl get namespace "$(K8S_NAMESPACE)" >/dev/null 2>&1 || kubectl create namespace "$(K8S_NAMESPACE)"

.PHONY: ns-down
ns-down:
	@kubectl delete namespace "$(K8S_NAMESPACE)" --ignore-not-found

.PHONY: deps-install
deps-install: ns-up
	@helm repo add bitnami https://charts.bitnami.com/bitnami >/dev/null 2>&1 || true
	@helm repo add nats https://nats-io.github.io/k8s/helm/charts >/dev/null 2>&1 || true
	@helm repo update >/dev/null
	@helm upgrade --install "$(POSTGRES_RELEASE)" bitnami/postgresql \
		--namespace "$(K8S_NAMESPACE)" \
		--values "$(POSTGRES_VALUES_FILE)"
	@helm upgrade --install "$(NATS_RELEASE)" nats/nats \
		--namespace "$(K8S_NAMESPACE)" \
		--values "$(NATS_VALUES_FILE)"

.PHONY: deps-uninstall
deps-uninstall:
	@helm uninstall "$(POSTGRES_RELEASE)" --namespace "$(K8S_NAMESPACE)" >/dev/null 2>&1 || true
	@helm uninstall "$(NATS_RELEASE)" --namespace "$(K8S_NAMESPACE)" >/dev/null 2>&1 || true

.PHONY: image-load
image-load:
	@if [[ -z "$(SERVICE)" ]]; then \
		echo "Usage: make image-load SERVICE=<service>"; \
		exit 1; \
	fi
	@if [[ ! -d "$(REPO_ROOT)/services/$(SERVICE)" ]]; then \
		echo "Unknown service: $(SERVICE)"; \
		exit 1; \
	fi
	@$(MAKE) -C "services/$(SERVICE)" containerise IMAGE_REPO="$(IMAGE_REPO)" IMAGE_TAG="$(IMAGE_TAG)"
	@k3d image import "$(IMAGE_REPO)/$(SERVICE)-service:$(IMAGE_TAG)" -c "$(K3D_CLUSTER_NAME)"

.PHONY: deploy
deploy:
	@if [[ -z "$(SERVICE)" ]]; then \
		echo "Usage: make deploy SERVICE=<service>"; \
		exit 1; \
	fi
	@if [[ ! -d "$(HELM_CHARTS_DIR)/$(SERVICE)" ]]; then \
		echo "Helm chart not found for service: $(SERVICE)"; \
		exit 1; \
	fi
	@helm upgrade --install "$(SERVICE)" "$(HELM_CHARTS_DIR)/$(SERVICE)" \
		--namespace "$(K8S_NAMESPACE)" \
		--create-namespace \
		--values "$(HELM_CHARTS_DIR)/$(SERVICE)/values.yaml" \
		--values "$(HELM_CHARTS_DIR)/$(SERVICE)/values-local.yaml"

.PHONY: deploy-all
deploy-all:
	@for s in $(DEPLOYABLE_SERVICES); do \
		$(MAKE) image-load SERVICE=$$s; \
		$(MAKE) deploy SERVICE=$$s; \
	done

.PHONY: wait-deps
wait-deps:
	@echo "Waiting for dependency pods to become Ready..."
	@kubectl wait --namespace "$(K8S_NAMESPACE)" \
		--for=condition=Ready pod \
		--selector "app.kubernetes.io/instance=$(POSTGRES_RELEASE)" \
		--timeout="$(WAIT_TIMEOUT)"
	@kubectl wait --namespace "$(K8S_NAMESPACE)" \
		--for=condition=Ready pod \
		--selector "app.kubernetes.io/instance=$(NATS_RELEASE)" \
		--timeout="$(WAIT_TIMEOUT)"

.PHONY: wait-services
wait-services:
	@for s in $(DEPLOYABLE_SERVICES); do \
		echo "Waiting for deployment/$$s rollout..."; \
		kubectl rollout status "deployment/$$s" --namespace "$(K8S_NAMESPACE)" --timeout="$(WAIT_TIMEOUT)"; \
	done

.PHONY: wait-ingress
wait-ingress:
	@echo "Waiting for ingress endpoint $(HEALTHCHECK_URL)..."
	@attempts=$$(( $(patsubst %s,%,$(WAIT_TIMEOUT)) / $(WAIT_RETRY_SECONDS) )); \
	if [[ $$attempts -lt 1 ]]; then attempts=1; fi; \
	last_code="000"; \
	for i in $$(seq 1 $$attempts); do \
		http_code=$$(curl --silent --show-error --output /dev/null --write-out "%{http_code}" "$(HEALTHCHECK_URL)" || true); \
		last_code="$$http_code"; \
		if [[ "$$http_code" =~ ^2[0-9][0-9]$$ ]]; then \
			echo "Ingress endpoint is reachable (HTTP $$http_code)."; \
			exit 0; \
		fi; \
		sleep "$(WAIT_RETRY_SECONDS)"; \
	done; \
	echo "Timed out waiting for ingress endpoint $(HEALTHCHECK_URL) (last code: $$last_code)"; \
	exit 1

.PHONY: stack-up
stack-up: cluster-up ns-up deps-install wait-deps deploy-all wait-services wait-ingress

.PHONY: stack-down
stack-down:
	@for s in $(DEPLOYABLE_SERVICES); do \
		helm uninstall $$s --namespace "$(K8S_NAMESPACE)" >/dev/null 2>&1 || true; \
	done
	@$(MAKE) deps-uninstall
