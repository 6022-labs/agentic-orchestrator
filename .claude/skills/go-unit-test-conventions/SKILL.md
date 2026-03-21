---
name: go-unit-test-conventions
description: Enforce consistent Go unit test conventions for this repository, including when_*_test.go naming, Given/Should nested subtests, BeforeEach suite setup, and go.uber.org/mock/gomock expectations. Use when creating or refactoring tests in tests/*_unit_tests and when generating mocks with deepsearch-mockgen.
---

# Go Unit Test Conventions

Apply this workflow to keep unit tests structurally consistent across the repository.

## Workflow

1. Identify test type and pick a template from `references/unit-test-templates.md`:
   - use case or service behavior test
   - request `Validate()` test
   - settings loader test
2. Keep naming aligned with existing conventions:
   - file: `when_<behavior>_test.go`
   - top test function: `TestWhen<Behavior>`
   - suite type: `When<Behavior>TestingSuite`
   - setup helper: `When<Behavior>BeforeEach`
3. Use `go.uber.org/mock/gomock` for mock creation and expectations in suite-based tests.
4. Structure scenarios as nested `t.Run` blocks with `Given ...` then `Should ...`, and call `t.Parallel()` at each level used in existing tests.
5. Keep setup localized:
   - `BeforeEach` builds SUT and all mocks.
   - scenario-specific expectations stay in `initSuite` closures.
   - default `initSuite` signature is `func(suite *When<Behavior>TestingSuite) ...`.
   - generate scenario-only random values inside `initSuite`.
   - pass extra parameters to `initSuite` only when those values are intentionally created outside and used in `Should` assertions/calls.
   - return only values consumed by `Should`; never return placeholder values that are ignored.
6. Reuse existing test utilities whenever possible:
   - `tests/common_tests`
   - `tests/common_tests/model_builders`
7. Regenerate mocks when interfaces change:
   - `deepsearch-mockgen -S ./src -O ./tests -A -P`
   - CLI repository: https://github.com/FournyP/deepsearch-mockgen-cli
8. Run targeted unit tests for touched packages before finishing:
   - `go test ./tests/<package>...`
   - Agents must rerun relevant tests after each fix until they pass.

## Required Structure Rules

- Keep package names external (`<pkg>_test`) when that is the local convention.
- Prefer explicit assertions with `testify/assert` (`assert.NoError`, `assert.ErrorIs`, `assert.Equal`, `assert.Len`, `assert.Nil`).
- Use `gomock.Any()` only when argument precision is not the behavior under test.
- Use `DoAndReturn` for argument inspection when behavior requires verification of transformed values.
- Keep one behavioral story per `Given` block; avoid mixing unrelated assertions in one scenario.
- Keep `Given` wording domain-centric (business state/behavior), not infrastructure-centric.
- Prefer model-centered `Given` wording whenever possible: start with the model under test (for example `Given model does not exists`, `Given saving model fails`).
- Instantiate domain models with builders from `tests/common_tests/model_builders` whenever a builder exists.
- Do not instantiate models directly with `&models.<Type>{}` or `models.<Type>{}` when a builder exists.

## Mocking Rules

- Import `go.uber.org/mock/gomock`.
- Instantiate one controller per test suite invocation: `gomock.NewController(t)`.
- Build mocks through generated constructors: `NewMock<Interface>(mockController)`.
- Define expectations with `.EXPECT()` and explicit `.Times(1)` when exact call count matters.

## Template Reference

Use `references/unit-test-templates.md` as the default source for copy-paste scaffolds.
