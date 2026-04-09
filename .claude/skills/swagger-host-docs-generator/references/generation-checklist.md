# Generation Checklist

Use this checklist after regenerating host docs.

1. Confirm each targeted host completed all generation steps without errors.
2. Confirm each host `docs/` directory contains:
   - `docs.go`
   - `swagger.json`
   - `swagger.yaml`
   - `swagger.v3.json`
   - `swagger.v3.yaml`
3. Confirm route/handler annotation changes are reflected in generated files.
4. Confirm no temporary `.tmp-v3` directory remains under host `docs/`.
5. Review `git diff` for generated artifacts only; avoid unrelated edits.
