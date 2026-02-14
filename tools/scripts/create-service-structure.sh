#!/usr/bin/env bash
set -euo pipefail

# create-service-structure.sh
#
# Usage:
#   tools/scripts/create-service-structure.sh identity
#
# Creates a new service folder structure under: services/<service>/
# Each service is a fully independent Go module:
#   module github.com/woffVienna/proteon-cursor/services/<service>

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <service-name>"
  exit 1
fi

SERVICE="$1"

# --- locate repo root (assumes script lives in tools/scripts/) ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

MODULE_PREFIX="github.com/woffVienna/proteon-cursor"
ROOT="${REPO_ROOT}/services/${SERVICE}"

if [[ -e "$ROOT" ]]; then
  echo "Error: '${ROOT}' already exists."
  exit 1
fi

# --- folders (folders first) ---
mkdir -p "${ROOT}/api/.generated"
mkdir -p "${ROOT}/api/events"
mkdir -p "${ROOT}/api/schemas"

mkdir -p "${ROOT}/cmd/${SERVICE}"

mkdir -p "${ROOT}/internal/adapters/db/migrations"
mkdir -p "${ROOT}/internal/adapters/db/postgres"

mkdir -p "${ROOT}/internal/adapters/http/generated"
mkdir -p "${ROOT}/internal/adapters/http/handlers"
mkdir -p "${ROOT}/internal/adapters/http/mapping"
mkdir -p "${ROOT}/internal/adapters/http/middleware"

mkdir -p "${ROOT}/internal/adapters/nats/consumer"
mkdir -p "${ROOT}/internal/adapters/nats/mapping"
mkdir -p "${ROOT}/internal/adapters/nats/publisher"

mkdir -p "${ROOT}/internal/application/dto"
mkdir -p "${ROOT}/internal/application/interfaces"
mkdir -p "${ROOT}/internal/application/services"

mkdir -p "${ROOT}/internal/domain/model"
mkdir -p "${ROOT}/internal/domain/rules"

mkdir -p "${ROOT}/internal/platform/buildinfo"
mkdir -p "${ROOT}/internal/platform/health"
mkdir -p "${ROOT}/internal/platform/shutdown"

mkdir -p "${ROOT}/test/contract"
mkdir -p "${ROOT}/test/integration"

# --- files (files after folders) ---
touch "${ROOT}/api/oapi-codegen.yaml"
touch "${ROOT}/api/openapi.yml"

# Keep empty directories in git (optional but useful)
touch "${ROOT}/api/events/.keep"
touch "${ROOT}/api/schemas/.keep"
touch "${ROOT}/internal/domain/model/.keep"
touch "${ROOT}/internal/domain/rules/.keep"
touch "${ROOT}/internal/platform/buildinfo/.keep"
touch "${ROOT}/internal/platform/health/.keep"
touch "${ROOT}/internal/platform/shutdown/.keep"
touch "${ROOT}/test/contract/.keep"
touch "${ROOT}/test/integration/.keep"

# Minimal, compilable main.go
cat > "${ROOT}/cmd/${SERVICE}/main.go" <<EOF
package main

import "log"

func main() {
	log.Println("${SERVICE} service starting...")
}
EOF

touch "${ROOT}/internal/adapters/db/db.go"
touch "${ROOT}/internal/adapters/http/server.go"
touch "${ROOT}/internal/adapters/nats/nats.go"

touch "${ROOT}/internal/domain/errors.go"
touch "${ROOT}/internal/domain/events.go"

touch "${ROOT}/README.md"

# Initialize a real go.mod (go.sum will be generated later by 'go mod tidy')
(
  cd "${ROOT}"
  go mod init "${MODULE_PREFIX}/services/${SERVICE}"
)

echo "Created service structure at: ${ROOT}"
echo "Initialized Go module: ${MODULE_PREFIX}/services/${SERVICE}"
echo "Next: (cd ${ROOT} && go mod tidy) after adding dependencies."
