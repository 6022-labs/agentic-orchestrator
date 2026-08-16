---
name: go-http-client-test-conventions
description: Enforce consistent Go HTTP client unit test conventions for HAND-WRITTEN clients only (e.g. conversation_broker_http_remote_agent, slack_bridge_http_slack) using net/http/httptest fake servers, when_*_test.go naming, Given/Should nested subtests, and BeforeEach/AfterEach suite setup. Do NOT use for generated oapi-codegen client wrappers — those are tested with gomock on clients.ClientWithResponsesInterface, never httptest.
---

# Go HTTP Client Test Conventions

**Scope: hand-written HTTP clients only** (today: `conversation_broker_http_remote_agent`, `slack_bridge_http_slack`). Adapters built on generated `oapi_client.gen.go` clients test their wrapper services with a gomock `MockClientWithResponsesInterface` instead — no httptest, no transport tests (the transport is generated code).

Apply this workflow to keep hand-written HTTP client unit tests structurally consistent across the repository. These tests exercise a thin client wrapper against a real `httptest.Server` so serialization, URL building, method, headers, and status-code handling are verified end-to-end at the transport boundary — without mocking `http.RoundTripper`.

## Workflow

1. Pick the template from `references/http-client-test-templates.md` that matches what the client under test does (request + decoded response, or fire-and-forget request).
2. Keep naming aligned with existing conventions:
   - file: `when_<behavior>_test.go` (e.g. `when_calling_completions_test.go`)
   - top test function: `TestWhen<Behavior>`
   - suite type: `When<Behavior>TestingSuite`
   - setup helper: `When<Behavior>BeforeEach`
   - teardown helper: `When<Behavior>AfterEach`
3. Stand up a stub upstream with `httptest.NewServer` in `BeforeEach`. The handler must delegate to a `suite.handler http.HandlerFunc` so each scenario can install its own behavior; fall back to `http.NotFound` when no handler is set so missing setup fails loudly.
4. Construct the SUT with the real `http.DefaultClient` (or equivalent) and point it at `suite.server.URL`. Inject a no-op logger (`zap.NewNop()`) to keep test output clean.
5. Close the server in `AfterEach` via `defer WhenCallingXAfterEach(suite)` to avoid leaked listeners.
6. Structure scenarios as nested `t.Run` blocks with `Given ...` then `Should ...`, and call `t.Parallel()` at each level that other tests in the repo parallelize.
7. Keep setup localized:
   - `BeforeEach` builds SUT, server, and wiring.
   - scenario-specific behavior (status, body, latency) lives in an `initSuite` closure that assigns `suite.handler`.
   - capture request facts (`method`, `path`, `body`, headers) into scenario-scoped variables declared inside the `Given` block, not on the suite.
8. Run targeted unit tests for touched packages before finishing:
   - `go test ./tests/<package>...`

## Required Structure Rules

- Keep package names external (`<pkg>_test`) when that is the local convention.
- Use `net/http/httptest.NewServer` — do not hand-roll a `net.Listener` or mock `http.RoundTripper` for this class of test.
- Route every request through a single indirection on the suite (`suite.handler`) so scenarios can swap behavior without rebuilding the server.
- Prefer explicit assertions with `testify/assert` (`assert.NoError`, `assert.Error`, `assert.Equal`, `assert.Contains`).
- Assert the observable HTTP contract, not the client internals:
  - method (`r.Method`)
  - path and query (`r.URL.Path`, `r.URL.RawQuery`)
  - relevant headers (`r.Header.Get(...)`)
  - request body — decode with `json.Unmarshal` when asserting structure, use `assert.Contains(string(body), ...)` only for quick substring checks.
- One behavioral story per `Given` block. Separate `Given upstream returns 2xx` from `Given upstream returns non-2xx`; do not fold error cases into the happy-path scenario.
- Cover at minimum: happy path (decoded response + captured request) and non-2xx (returns error). Add `Given upstream returns malformed body` when the client decodes a response.

## Stub Server Rules

- Always write the `Content-Type` header before encoding a JSON body (`w.Header().Set("Content-Type", "application/json")`).
- Use the correct status helpers: `w.WriteHeader(http.StatusNoContent)` for empty success, `http.StatusInternalServerError` for failure scenarios.
- Read the request body once with `io.ReadAll(r.Body)`; capture into a scenario-local `[]byte` so assertions can run after the handler returns.
- Never rely on goroutine-unsafe state shared across parallel subtests — keep captured variables inside the `Given` closure, not on the suite, unless the suite itself is built fresh per `Should` (which is the convention here).

## Template Reference

Use `references/http-client-test-templates.md` as the default source for copy-paste scaffolds.
