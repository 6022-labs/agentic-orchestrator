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

for DB in mariadb mysql postgres; do
    DIR="$MIGRATIONS_ROOT/$DB"
    echo "Generating migration for $DB..."
    atlas migrate diff "$MIGRATION_NAME" --env gorm-"$DB"
    echo "Generated migration for $DB in directory: $DIR"
    echo
done

TIMESTAMP=$(date +"%Y%m%d%H%M%S")

for DB in mariadb mysql postgres sqlserver; do
    DIR="$MIGRATIONS_ROOT/$DB"

    LATEST_FILE=$(ls -1t "$DIR"/*.sql 2>/dev/null | head -n 1)
    if [ -z "$LATEST_FILE" ]; then
        echo "No SQL file found in $DIR, skipping."
        continue
    fi

    NAME_PART=$(basename "$LATEST_FILE" | cut -d'_' -f2-)
    NEW_NAME="${TIMESTAMP}_${NAME_PART}"
    NEW_PATH="$DIR/$NEW_NAME"
    mv "$LATEST_FILE" "$NEW_PATH"
    echo "Renamed $LATEST_FILE to $NEW_PATH"

    atlas migrate hash --dir file://"$DB"
done