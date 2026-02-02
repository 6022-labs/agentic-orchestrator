#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

find "$repo_root" -name go.mod -not -path "*/vendor/*" | while IFS= read -r modfile; do
  service_dir="$(dirname "$modfile")"
  echo "Running go mod tidy in $service_dir"
  (cd "$service_dir" && go mod tidy)
done
