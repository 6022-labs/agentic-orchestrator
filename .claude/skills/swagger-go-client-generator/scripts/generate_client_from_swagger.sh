#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  generate_client_from_swagger.sh \
    [--spec-path <path>] \
    [--output-dir <dir>] \
    [--output-file <filename>] \
    [--client-package <package>] \
    [--generate <types,client>] \
    [--oapi-codegen-version <version>]
EOF
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

write_generation_readme() {
  local target_dir="$1"
  local readme_path="${target_dir}/README.md"

  cat > "${readme_path}" <<'EOF'
# Generated package

This package was generated using the `swagger-go-client-generator` Copilot skill.
The generated code is based on a local OpenAPI v3 spec and `oapi-codegen`.

Do not edit generated files manually; regenerate from the source OpenAPI spec instead.
EOF
}

SPEC_PATH="src/copy_liquidity_providing_synchronizer_host/docs/swagger.v3.yaml"
OUTPUT_DIR="src/copy_liquidity_providing_backoffice_http_synchronizer/clients"
OUTPUT_FILE="oapi_client.gen.go"
CLIENT_PACKAGE="clients"
GENERATE_TARGETS="types,client"
OAPI_CODEGEN_VERSION="v2.5.1"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --spec-path)
      SPEC_PATH="${2:-}"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="${2:-}"
      shift 2
      ;;
    --output-file)
      OUTPUT_FILE="${2:-}"
      shift 2
      ;;
    --client-package)
      CLIENT_PACKAGE="${2:-}"
      shift 2
      ;;
    --generate)
      GENERATE_TARGETS="${2:-}"
      shift 2
      ;;
    --oapi-codegen-version)
      OAPI_CODEGEN_VERSION="${2:-}"
      shift 2
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

if [[ -z "$SPEC_PATH" || -z "$OUTPUT_DIR" || -z "$OUTPUT_FILE" || -z "$CLIENT_PACKAGE" || -z "$GENERATE_TARGETS" ]]; then
  echo "missing required values" >&2
  exit 1
fi

require_cmd go

if [[ ! -f "${SPEC_PATH}" ]]; then
  echo "OpenAPI spec not found: ${SPEC_PATH}" >&2
  exit 1
fi

mkdir -p "$OUTPUT_DIR"

echo "[1/1] Generating Go client with oapi-codegen"
go run "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}" \
  -generate "${GENERATE_TARGETS}" \
  -package "${CLIENT_PACKAGE}" \
  -o "${OUTPUT_DIR}/${OUTPUT_FILE}" \
  "${SPEC_PATH}"

write_generation_readme "${OUTPUT_DIR}"

echo "Client generation completed"
