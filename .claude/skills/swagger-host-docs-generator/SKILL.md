---
name: swagger-host-docs-generator
description: Generate and refresh Swagger/OpenAPI documentation for Go host services that expose a docs generation script. Use when asked to regenerate host API docs, produce swagger.v3 files, or update docs artifacts under *_host/docs.
---

# Swagger Host Docs Generator

Generate host Swagger/OpenAPI documentation by running each host's local `docs/generate_docs.sh` script.

## Workflow

1. Identify target hosts:
   - Use all discovered hosts by default.
   - Or restrict to specific hosts (example: `gateway_host`).
2. Run the bundled script from repository root:
   - `./scripts/generate_swagger_docs_for_hosts.sh`
   - `./scripts/generate_swagger_docs_for_hosts.sh --host gateway_host`
3. Verify generated artifacts in each host `docs/` directory:
   - `docs.go`
   - `swagger.json`
   - `swagger.yaml`
   - `swagger.v3.json`
   - `swagger.v3.yaml`
4. Review and commit changed docs artifacts.

## Commands

```bash
# Regenerate docs for every discovered host
./scripts/generate_swagger_docs_for_hosts.sh

# Regenerate docs for a specific host
./scripts/generate_swagger_docs_for_hosts.sh --host gateway_host

# Preview without running generation
./scripts/generate_swagger_docs_for_hosts.sh --dry-run
```

## Notes

- Keep generation host-local by calling each host's `docs/generate_docs.sh` instead of hand-written `swag` commands.
- Run from repository root unless passing `--repo-root`.
- Use `references/generation-checklist.md` for verification and troubleshooting.
