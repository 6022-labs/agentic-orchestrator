#!/usr/bin/env sh
set -eu

# Regenerates the bridge HTTP client from slack_bridge_internal_host's OpenAPI v3 spec.
# Run channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs/generate_docs.sh first if the API changed.

CLIENTS_DIR="$(CDPATH= cd -- "$(dirname "$0")" && pwd)"
SPEC_PATH="${CLIENTS_DIR}/../../../../channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs/swagger.v3.yaml"
OUTPUT_FILE="${CLIENTS_DIR}/oapi_client.gen.go"

OAPI_CODEGEN_VERSION="v2.5.1"

echo "Generating bridge callback client with oapi-codegen ${OAPI_CODEGEN_VERSION}"
(
	cd "${CLIENTS_DIR}"
	go run "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}" \
		-generate types,client \
		-package clients \
		-o "${OUTPUT_FILE}" \
		"${SPEC_PATH}"
)

echo "Done"
