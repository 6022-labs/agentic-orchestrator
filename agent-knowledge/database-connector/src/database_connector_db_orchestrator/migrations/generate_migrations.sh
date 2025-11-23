#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 1 ]; then
    echo "Usage: $0 <migration-name>" >&2
    exit 1
fi

MIGRATION_NAME="$1"
MIGRATIONS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if ! command -v atlas >/dev/null 2>&1; then
    echo "Error: atlas CLI not found in PATH" >&2
    exit 1
fi

atlas migrate diff "$MIGRATION_NAME" --env gorm