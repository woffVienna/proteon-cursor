#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${CLUSTER_NAME:-proteon}"
AGENTS="${AGENTS:-2}"

cluster_exists() {
  k3d cluster list --no-headers 2>/dev/null | awk '{print $1}' | grep -Fxq "${CLUSTER_NAME}"
}

create_cluster() {
  if cluster_exists; then
    echo "k3d cluster '${CLUSTER_NAME}' already exists."
    return 0
  fi

  echo "Creating k3d cluster '${CLUSTER_NAME}'..."
  k3d cluster create "${CLUSTER_NAME}" \
    --agents "${AGENTS}" \
    --port "8080:80@loadbalancer" \
    --port "8443:443@loadbalancer"
}

delete_cluster() {
  if ! cluster_exists; then
    echo "k3d cluster '${CLUSTER_NAME}' does not exist."
    return 0
  fi

  echo "Deleting k3d cluster '${CLUSTER_NAME}'..."
  k3d cluster delete "${CLUSTER_NAME}"
}

usage() {
  cat <<EOF
Usage: $0 <up|down>

Commands:
  up    Create the local k3d cluster if missing
  down  Delete the local k3d cluster if present
EOF
}

main() {
  local cmd="${1:-}"
  case "${cmd}" in
    up)
      create_cluster
      ;;
    down)
      delete_cluster
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
