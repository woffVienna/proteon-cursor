#!/usr/bin/env bash
set -Eeuo pipefail

trap 'echo "ERROR: Script failed at line $LINENO. Command: $BASH_COMMAND" >&2' ERR

SOURCE_DIR="$(pwd)"
PARENT_DIR="$(dirname "$SOURCE_DIR")"
REPO_NAME="$(basename "$SOURCE_DIR")"
TIMESTAMP="$(date +"%Y%m%d_%H%M%S")"

SCRIPT_NAME="$(basename "$0")"

CLEAN_DIR="${PARENT_DIR}/${REPO_NAME}_clean_${TIMESTAMP}"
ZIP_FILE="${PARENT_DIR}/${REPO_NAME}_clean_${TIMESTAMP}.zip"

echo "== Package monorepo for AI =="
echo "Source:     $SOURCE_DIR"
echo "Parent:     $PARENT_DIR"
echo "Clean dir:  $CLEAN_DIR"
echo "Zip file:   $ZIP_FILE"
echo

# Dependency checks
for cmd in rsync zip; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: Required command '$cmd' not found in PATH." >&2
    exit 1
  fi
done

echo "Step 1/3: Creating clean copy..."
mkdir -p "$CLEAN_DIR"

rsync -a --delete \
  --exclude=".git" \
  --exclude=".gitignore" \
  --exclude=".idea" \
  --exclude=".vscode" \
  --exclude="node_modules" \
  --exclude="vendor" \
  --exclude="bin" \
  --exclude=".build" \
  --exclude="dist" \
  --exclude="coverage" \
  --exclude="tmp" \
  --exclude="*.log" \
  --exclude="*.out" \
  --exclude="*.env" \
  --exclude="*.local" \
  --exclude="**/.DS_Store" \
  --exclude="$SCRIPT_NAME" \
  "$SOURCE_DIR/" "$CLEAN_DIR/"

echo "Step 2/3: Creating zip..."
rm -f "$ZIP_FILE"
(
  cd "$PARENT_DIR"
  zip -rq "$(basename "$ZIP_FILE")" "$(basename "$CLEAN_DIR")"
)

echo "Step 3/3: Removing temporary clean directory..."
rm -rf "$CLEAN_DIR"

echo
echo "DONE ✅"
echo "Created: $ZIP_FILE"
ls -lh "$ZIP_FILE"