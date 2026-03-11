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
DEPLOYABLE_SERVICES := identity api-gateway auth backoffice-gateway

IMAGE_REPO ?= proteon
IMAGE_TAG ?= dev
WAIT_TIMEOUT ?= 180s
WAIT_RETRY_SECONDS ?= 2
HEALTHCHECK_URL ?= http://localhost:8080/v1/health

.PHONY: help-infra
help-infra:
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
