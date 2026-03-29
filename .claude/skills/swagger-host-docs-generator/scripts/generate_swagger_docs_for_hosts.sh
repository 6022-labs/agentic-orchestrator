#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  generate_swagger_docs_for_hosts.sh [--repo-root <path>] [--host <host-dir>]... [--dry-run]

Examples:
  generate_swagger_docs_for_hosts.sh
  generate_swagger_docs_for_hosts.sh --host collateral_dashboard_host
  generate_swagger_docs_for_hosts.sh --repo-root /path/to/repo --host collateral_image_generator_host
USAGE
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

matches_host_filter() {
  local host_name="$1"
  if [[ "${#HOST_FILTERS[@]}" -eq 0 ]]; then
    return 0
  fi

  local filter
  for filter in "${HOST_FILTERS[@]}"; do
    if [[ "$host_name" == "$filter" ]]; then
      return 0
    fi
  done

  return 1
}

REPO_ROOT="$(pwd)"
DRY_RUN=0
HOST_FILTERS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo-root)
      REPO_ROOT="${2:-}"
      shift 2
      ;;
    --host)
      HOST_FILTERS+=("${2:-}")
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$REPO_ROOT" ]]; then
  echo "repo root cannot be empty" >&2
  exit 1
fi

require_cmd find
require_cmd sort
require_cmd go

if [[ ! -d "$REPO_ROOT/src" ]]; then
  echo "src directory not found under repo root: $REPO_ROOT" >&2
  exit 1
fi

ALL_SCRIPTS=()
while IFS= read -r script_path; do
  ALL_SCRIPTS+=("$script_path")
done < <(find "$REPO_ROOT/src" -type f -path '*_host/docs/generate_docs.sh' | sort)

if [[ "${#ALL_SCRIPTS[@]}" -eq 0 ]]; then
  echo "no host docs generation scripts found under $REPO_ROOT/src" >&2
  exit 1
fi

TARGET_SCRIPTS=()
for script_path in "${ALL_SCRIPTS[@]}"; do
  host_dir="$(basename "$(dirname "$(dirname "$script_path")")")"
  if matches_host_filter "$host_dir"; then
    TARGET_SCRIPTS+=("$script_path")
  fi
done

if [[ "${#TARGET_SCRIPTS[@]}" -eq 0 ]]; then
  echo "no matching host docs generation scripts found" >&2
  echo "requested hosts: ${HOST_FILTERS[*]:-<none>}" >&2
  exit 1
fi

echo "Found ${#TARGET_SCRIPTS[@]} host docs generation script(s)."

index=1
for script_path in "${TARGET_SCRIPTS[@]}"; do
  host_dir="$(basename "$(dirname "$(dirname "$script_path")")")"
  echo "[$index/${#TARGET_SCRIPTS[@]}] ${host_dir}"

  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "  dry-run: sh $script_path"
  else
    sh "$script_path"
  fi

  index=$((index + 1))
done

echo "Swagger/OpenAPI generation completed."
