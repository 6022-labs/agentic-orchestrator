---
name: swagger-go-client-generator
description: Regenerate the generated oapi-codegen Go HTTP clients (`src/*_http_*/clients/oapi_client.gen.go`) from the host OpenAPI v3 specs in this repository. Use when a host API contract changes and one or more generated clients must be refreshed, or when adding a new generated-client adapter.
---

# Swagger Go Client Generator

Every generated-client HTTP adapter owns a `clients/generate_client.sh` next to its
`oapi_client.gen.go`. Regeneration is always: refresh the provider host's spec, then
run the consumer adapters' scripts, then regenerate mocks.

## Adapters and their provider hosts

| Adapter (clients/generate_client.sh) | Provider host spec |
|---|---|
| `gateway/src/gateway_http_agent_identity_registry` | `agent-identity/src/registry_host` |
| `agent-swarm/src/conversation_broker_http_agent_identity_registry` | `agent-identity/src/registry_host` |
| `channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry` | `agent-identity/src/registry_host` |
| `gateway/src/gateway_http_agent_swarm_conversation_broker` | `agent-swarm/src/conversation_broker_host` |
| `channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker` | `agent-swarm/src/conversation_broker_host` |
| `gateway/src/gateway_http_channels_and_integrations_bridge` | `channels-and-integrations/slack-bridge/src/slack_bridge_internal_host` |
| `agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge` | `channels-and-integrations/slack-bridge/src/slack_bridge_internal_host` |

## Workflow (from repo root)

1. If the provider's endpoints changed, regenerate its OpenAPI spec first:
   - `sh <provider-host>/docs/generate_docs.sh` (or `bash scripts/generate_swagger_docs_for_hosts.sh --host <host_name>`)
2. Regenerate every consumer client of that host:
   - `sh <adapter>/clients/generate_client.sh`
3. Regenerate the adapter's mock (from the adapter's module root):
   - `go run go.uber.org/mock/mockgen@v0.6.0 -source=src/<adapter>/clients/oapi_client.gen.go -destination=tests/<adapter>_mocks/clients_mocks/client_with_responses_mock.go -package=clients_mocks`
4. Build + run the adapter's unit tests.

## Generation Guarantees

- Keep `oapi-codegen` pinned to `v2.5.1` and `-generate types,client`.
- Keep output at `<adapter>/clients/oapi_client.gen.go`, package `clients`.
- Wrapper services depend on `clients.ClientWithResponsesInterface`; request bodies
  go through `json.Marshal` + the `...WithBodyWithResponse` variants (swag emits
  `oneOf` bodies that oapi-codegen turns into setterless unions).
- Each adapter's `clients/README.md` documents its exact refresh commands.

## New adapter

Copy `gateway/src/gateway_http_agent_identity_registry/` as the template (clients/
mappings/ helpers/ services/ settings/ configurations/), adjust the SPEC_PATH depth
(slack-bridge adapters sit one directory deeper: five `..` instead of four), and add
`github.com/oapi-codegen/runtime` to the module if missing.
