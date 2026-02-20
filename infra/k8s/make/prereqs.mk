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

.PHONY: help-prereqs
help-prereqs:
	@echo "  check-prereqs - verify required local stack binaries are installed"
	@echo "  setup-darwin-arm64-prereqs - install local stack prerequisites on darwin-arm64"

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
