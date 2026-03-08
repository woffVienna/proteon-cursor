#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "ERROR: Script failed at line $LINENO. Command: $BASH_COMMAND" >&2' ERR

SOURCE_DIR="$(pwd)"
PARENT_DIR="$(dirname "$SOURCE_DIR")"
REPO_NAME="$(basename "$SOURCE_DIR")"
TIMESTAMP="$(date +"%Y%m%d_%H%M%S")"

ARCH_DIR="$SOURCE_DIR/docs/architecture"
OUTPUT_DIR="${PARENT_DIR}/${REPO_NAME}_architecture_${TIMESTAMP}"

echo "== Package architecture for AI =="
echo "Source:     $ARCH_DIR"
echo "Output:     $OUTPUT_DIR"
echo

if [ ! -d "$ARCH_DIR" ]; then
  echo "ERROR: docs/architecture/ not found in $SOURCE_DIR" >&2
  exit 1
fi

echo "Copying all .md files from docs/architecture/ into a flat folder..."
mkdir -p "$OUTPUT_DIR"

find "$ARCH_DIR" -name '*.md' -type f | while read -r file; do
  basename="$(basename "$file")"

  if [ -e "$OUTPUT_DIR/$basename" ]; then
    dir="$(basename "$(dirname "$file")")"
    basename="${dir}_${basename}"
  fi

  cp "$file" "$OUTPUT_DIR/$basename"
done

echo
echo "DONE"
echo "Created: $OUTPUT_DIR"
echo
ls -1 "$OUTPUT_DIR"
