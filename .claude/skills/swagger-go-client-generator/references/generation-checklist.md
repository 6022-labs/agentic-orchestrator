# Generation Checklist

Use this checklist when wiring client generation into a repository script.

1. Keep the source spec path explicit (`src/copy_liquidity_providing_synchronizer_host/docs/swagger.v3.yaml` by default).
2. Keep `oapi-codegen` pinned at `v2.5.1` unless generation output is revalidated.
3. Keep `-generate types,client` unless consumer code is updated for a different shape.
4. Keep package name `clients` to preserve import compatibility.
5. Keep output file explicit and deterministic (`oapi_client.gen.go`).
6. Add/update generated package `README.md` to indicate skill-based generation.
7. Regenerate docs (`go generate ./src/copy_liquidity_providing_synchronizer_host/docs`) before client generation when endpoints changed.
8. Regenerate the client whenever the OpenAPI schema changes.
