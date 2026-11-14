---
name: go-integration-test-conventions
description: Enforce consistent Go integration test conventions for this repository, including when_*_test.go naming, Given/Should nested subtests, BeforeEach/AfterEach setup around real infrastructure, and reuse of tests/common_tests/container_builders for PostgreSQL and Redis. Use when creating or refactoring tests in tests/*_integration_tests.
---

# Go Integration Test Conventions

Apply this workflow to keep integration tests consistent across repository packages.

## Workflow

1. Identify integration test category and start from `references/integration-test-templates.md`:
   - DB repository integration test (PostgreSQL + migrations)
   - Redis integration test (real Redis client + storage service)
   - factory/infrastructure integration test
2. Keep naming aligned with existing conventions:
   - file: `when_<behavior>_test.go`
   - top test function: `TestWhen<Behavior>`
   - suite type: `When<Behavior>TestingSuite`
   - setup helpers: `When<Behavior>BeforeEach` and `When<Behavior>AfterEach`
3. Use real infrastructure with shared helpers:
   - PostgreSQL DB: `container_builders.GetPostgresContainerBuilder().Create()` (`CreateTls()` for a TLS server) + `DropDatabase`
   - Redis: `container_builders.GetRedisContainerBuilder().Create()` + `DropConnection`
4. Build scenario setup in `BeforeEach`, always clean up with `defer ...AfterEach(...)`.
5. Keep scenario structure as nested `t.Run` blocks with `Given ...` then `Should ...`; call `t.Parallel()` at each level used by surrounding tests.
6. For DB integration tests, run migrations after creating the DB client:
   - resolve migrations path with `common_tests.FindMigrationFolderPath()`
   - call `dbClient.Migrate()`
7. Keep mocks minimal in integration tests:
   - allow mocks for cross-cutting adapters (for example metrics interfaces) with `AnyTimes()`
   - avoid mocking the primary dependency under test (database/redis behavior must be real)
8. Run targeted integration tests for touched packages:
   - `go test ./tests/<integration_package>...`
   - Agents must rerun relevant tests after each fix until they pass.

## Required Structure Rules

- Keep behavior-focused scenarios and assertions; avoid mixing unrelated flows in one `Given` block.
- Prefer `testify/assert` explicit checks (`assert.NoError`, `assert.ErrorIs`, `assert.Equal`, `assert.Len`, `assert.Nil`).
- Reuse `tests/common_tests/model_builders` for domain test data setup.
- Keep `Given` wording domain-centric (state/behavior) instead of repository/container implementation phrasing.
- Prefer model-centered `Given` wording whenever possible: start with the model under test (for example `Given model does not exists`, `Given saving model fails`).
- Build domain entities with builders from `tests/common_tests/model_builders` whenever a builder exists.
- Do not instantiate models directly with `&models.<Type>{}` or `models.<Type>{}` when a builder exists.
- Keep cleanup deterministic (`defer When<Behavior>AfterEach(suite)` inside each scenario).

## Infrastructure Rules

- Require Docker engine availability for testcontainers-backed integration tests.
- Keep test packages under `tests/*_integration_tests`.
- Keep container lifecycle in shared helpers (`tests/common_tests/container_builders`) rather than per-test ad hoc setup.

## Template Reference

Use `references/integration-test-templates.md` as the default source for copy-paste scaffolds.
