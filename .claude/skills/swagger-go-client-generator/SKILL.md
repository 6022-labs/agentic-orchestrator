---
name: swagger-go-client-generator
description: Regenerate the `copy_liquidity_providing_backoffice_http_synchronizer` Go HTTP client with `oapi-codegen` from the synchronizer OpenAPI v3 spec. Use when the synchronizer API contract changes and `src/copy_liquidity_providing_backoffice_http_synchronizer/clients/oapi_client.gen.go` must be refreshed.
---

# Swagger Go Client Generator

Regenerate deterministic Go clients with `oapi-codegen` from the local synchronizer OpenAPI v3 spec.

## Workflow

1. If endpoints changed, regenerate synchronizer OpenAPI docs first:
   - `go generate ./src/copy_liquidity_providing_synchronizer_host/docs`
2. Regenerate the backoffice synchronizer client:
   - `go generate ./src/copy_liquidity_providing_backoffice_http_synchronizer/clients`
3. Or run the bundled script for explicit inputs/overrides.
4. Re-run generation whenever the API contract changes.

## Command

```bash
./scripts/generate_client_from_swagger.sh \
  --spec-path src/copy_liquidity_providing_synchronizer_host/docs/swagger.v3.yaml \
  --output-dir src/copy_liquidity_providing_backoffice_http_synchronizer/clients
```

## Generation Guarantees

- Keep `oapi-codegen` pinned to `v2.5.1`.
- Keep `-generate types,client`.
- Keep output at `src/copy_liquidity_providing_backoffice_http_synchronizer/clients/oapi_client.gen.go`.
- Keep package name `clients` so existing imports remain compatible.

## Additional Guidance

Read `references/generation-checklist.md` when adapting this workflow to a repository-specific script.
