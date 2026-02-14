#!/usr/bin/env bash
set -euo pipefail

# go-work-use.sh
#
# Usage:
#   tools/scripts/go-work-use.sh identity
#
# Ensures:
#   - go.work exists at repo root
#   - libs/platform (if present) is added
#   - services/<service> is added

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <service-name>"
  exit 1
fi

SERVICE="$1"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

SERVICE_PATH="${REPO_ROOT}/services/${SERVICE}"
PLATFORM_PATH="${REPO_ROOT}/libs/platform"

if [[ ! -d "$SERVICE_PATH" ]]; then
  echo "Error: service '${SERVICE}' does not exist at ${SERVICE_PATH}"
  exit 1
fi

cd "$REPO_ROOT"

# Create go.work if it does not exist
if [[ ! -f "go.work" ]]; then
  echo "Initializing go.work..."
  go work init
fi

# Helper to check if a module is already in go.work
is_in_work() {
  go work edit -json | grep -q "\"DiskPath\": \"${1}\""
}

# Add libs/platform if present
if [[ -d "$PLATFORM_PATH" ]]; then
  if ! is_in_work "./libs/platform"; then
    echo "Adding libs/platform to go.work"
    go work use ./libs/platform
  fi
fi

# Add service module
if ! is_in_work "./services/${SERVICE}"; then
  echo "Adding services/${SERVICE} to go.work"
  go work use ./services/${SERVICE}
else
  echo "services/${SERVICE} already present in go.work"
fi

echo "go.work updated successfully."
