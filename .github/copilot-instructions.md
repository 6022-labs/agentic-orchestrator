# GitHub Copilot Instructions

## Priority Guidelines

1. **Version Compatibility**: Every service is pinned to Go 1.24.x (see `go.work` and the `go.mod` files in `gateway`, `agent-knowledge/database-connector`, `agent-identity/images-generator`, `agent-reasoning/llm-connector`, and `channels-and-integrations/slack-bridge`). Frontend code targets React 19 / Vite 6 / TypeScript 5.7.2 (`ui/package.json`) and the docs site targets Docusaurus 3.9.2 with TypeScript 5.6.2 (`docs/package.json`). Never introduce APIs beyond those versions.
2. **Context Files**: This document (`.github/copilot/copilot-instructions.md`) is the authoritative Copilot context. When it references patterns, use the cited source files (controllers under `*/mvc`, use cases under `*/use_cases`, repositories under `orchestrator_db/repositories`, UI components in `ui/src/components`, docs in `docs/docs`).
3. **Codebase Patterns**: Keep the existing layering: HTTP controllers → use cases → services/interfaces → repositories/clients for Go microservices, context/provider driven state + hooks + `react-hook-form` for the UI, and Markdown with Docusaurus front matter for docs. This mirrors the repo’s hexagonal architecture (adapters on the outside, domain/use-case logic in the center). When unsure, copy the most recent example from the same folder.
4. **Architectural Consistency**: Treat the repo as a microservices mono-repo that applies hexagonal boundaries. Preserve the host/driver split (`*_host`, `*_mvc`, `*_http_*`, `*_db`, `*_vault`, `*_redis`), keep adapters isolated from the core use cases, and never reach across a service boundary (e.g., gateway logic should call agent connectors via the HTTP client modules, not by importing their internals directly).
5. **Code Quality**: Balance maintainability, performance, security, accessibility, and testability the way current code does: structured logging with Zap, concurrency through goroutines + channels where already used, JWT + SIWE auth middleware, Radix/Tailwind UI primitives with accessible focus handling, and gomock/testify based Go tests.

## Technology Version Detection

- **Go toolchain**: `go.work` pins Go 1.24.3 for the workspace; slack-bridge currently uses 1.24.0. Use `go env GOVERSION` or the local `go.mod` before adopting features.
- **Gateway & agent services** (`gateway/go.mod`, `agent-knowledge/database-connector/go.mod`, `agent-reasoning/llm-connector/go.mod`, `agent-identity/images-generator/go.mod`, `channels-and-integrations/slack-bridge/go.mod`):
  - HTTP: `github.com/gofiber/fiber/v2 v2.52.8`, `github.com/gofiber/swagger v1.1.1`, `github.com/gofiber/contrib/fiberzap/v2 v2.1.6`.
  - Logging/DI/testing: `go.uber.org/zap v1.27.0`, `go.uber.org/dig v1.19.0`, `go.uber.org/mock v0.5.2`, `github.com/stretchr/testify v1.10.0`.
  - Data: `gorm.io/gorm v1.30.0`, `gorm.io/driver/postgres v1.6.0`, `github.com/pressly/goose/v3 v3.24.3`, `ariga.io/atlas-provider-gorm v0.5.3`.
  - Infra/security: `github.com/hashicorp/vault/api v1.20.0`, `github.com/redis/go-redis/v9 v9.10.0`, `github.com/golang-jwt/jwt/v5 v5.2.2`.
  - Blockchain: `github.com/ethereum/go-ethereum v1.15.11`, `github.com/wealdtech/go-ens/v3 v3.6.0`.
  - Testing env: `github.com/testcontainers/testcontainers-go v0.36.0` (some modules use 0.37.0).
- **UI (`ui/package.json`)**: React 19, React Router 7.2.0, Vite 6.2, TypeScript ~5.7.2, Tailwind 4.0.9 via `@tailwindcss/vite`, Radix UI 1.1.x, `@reown/appkit` 1.6.9, Ethers 6.13.5, Axios 1.8.1, `react-hook-form` 7.66.1, `zod` 4.1.13. Respect the `tsconfig.app.json` compiler flags (strict mode, bundler resolution) and the flat ESLint config in `ui/eslint.config.js`.
- **Docs (`docs/package.json`)**: `@docusaurus/core` 3.9.2, React 19, TypeScript ~5.6.2. Node 18+ is required (see `engines.node`).
- **Image generation** (`agent-identity/images-generator/go.mod`): `github.com/chromedp/chromedp v0.10.0` drives headless rendering; keep API usage compatible with that version.

## Context Files

- `.github/copilot/copilot-instructions.md` (this file) – canonical guidance.
- There are no other `.github/copilot` helpers yet; rely on the service-specific READMEs (e.g., `gateway/README.md`) and controller/use-case code for domain-specific nuances.

## Codebase Scanning Instructions

1. For Go changes, inspect the matching service folder (`gateway/src/...`, `agent-knowledge/database-connector/src/...`, `channels-and-integrations/slack-bridge/src/...`). Follow the established pipeline: controllers in `*_mvc`, request/response DTOs in `use_cases/requests|responses`, business logic in `use_cases`, adapters in `services`, persistence in `orchestrator_db/repositories`.
2. For front-end work, mirror the organization under `ui/src`: shared hooks/contexts live in `contexts/`, requests in `queries/`, typed models in `models/`, UI building blocks in `components/ui/`, and feature surfaces under `components/dashboard/**`.
3. For documentation, edit the MD/MDX files in `docs/docs/**`. Each doc begins with YAML front matter (`id`, `title`, `sidebar_label`, etc.) and is registered in `docs/sidebars.ts`.
4. Cross-service communication is handled through HTTP client packages (e.g., `gateway/src/gateway_http_agent_reasoning_compute`). Never bypass these clients; update or extend them instead.
5. When touching persistence, review existing migrations under `*/migrations` and the Atlas schema files before adding SQL or schema drift.

## Code Quality Standards

### Maintainability

- Keep controllers thin: follow `gateway/src/orchestrator_gateway_mvc/agent_controller.go` by only parsing inputs, logging via `zap.Logger`, and delegating to use cases.
- Preserve DTO sanitization patterns (`requests/PaginateRequest.Sanitize` and `lo.ToPtr` helpers) and reuse shared validation types such as `errors.ValidationError`.
- Prefer constructor helpers and builders (e.g., `tests/common_tests/model_builders/agent_builder.go`) for complex setups.

### Performance

- Match existing concurrency techniques: `use_cases/ListAgents` spins goroutines per agent and funnels results through typed channels—reuse that pattern whenever a use case fan-outs work.
- Lean on the existing cache abstraction (`gateway_redis/services/redis_storage.go`) before querying remote systems repeatedly.
- Favor streaming/limited logging helpers like `zapext.ArrayTruncateEnd` when dumping large payloads.

### Security

- Reuse the layered auth middleware stack: `middlewares.JWTAuthenticationInjectionMiddleware` populates `authenticated-wallet`, and `AuthenticationMandatoryMiddleware` enforces it. New endpoints must live behind the same guards unless they are explicitly public.
- Store and read secrets with `VaultSecretStorage` (`gateway_vault/services/vault_secret_storage.go`) instead of environment variables or config files.
- Issue parameterized DB queries through GORM repositories and keep blockchain signing via the existing services (`services/auth/signature/ECDSASigner`, `services/auth/jwt/JwtManager`).

### Accessibility

- In React, always leverage the primitives in `ui/src/components/ui` (`Button`, `Select`, `Form`, etc.) so that focus rings, ARIA attributes, and keyboard support remain consistent (see `Button` and its CVA variants in `components/ui/button.tsx`).
- Reuse semantic HTML as shown in `ui/src/pages/login.tsx` and cascade utility classes through the `cn` helper rather than hard-coding inline styles.

### Testability

- Follow the structure in `gateway/tests/orchestrator_gateway_unit_tests/**`: table-driven subtests with `t.Parallel()`, gomock controllers created via `gomock.NewController(t, gomock.WithOverridableExpectations())`, assertions with `testify/assert`, and fixture builders from `tests/common_tests`.
- Integration tests spin ephemeral infra via `tests/common_tests/container_builders` and `testcontainers-go`. New integration suites should set up and tear down containers the same way (see `gateway/tests/orchestrator_db_integration_tests/factories/when_creating_orchestrator_db_client_test.go`).

## Documentation Requirements

- Go HTTP handlers are documented with Swagger blocks (`@Summary`, `@Description`, `@Param`, `@Router`, etc.) directly above the method (see `agent_controller.go`). Match that tone and include every parameter and status code.
- Docusaurus docs start with YAML front matter and use Markdown/MDX plus Mermaid diagrams (`docs/docusaurus.config.ts` enables mermaid). Mirror the style in `docs/docs/developers/intro.md` when adding sections.
- Service-specific READMEs should be updated whenever you change their bootstrapping commands (`docker-compose`, `.env` requirements, etc.).

## Testing Approach

### Unit Testing

- Co-locate Go unit tests under `gateway/tests/<domain>_unit_tests`. Each feature uses a `When<X>TestSuite` helper and arranges mocks via `*mocks` packages generated by gomock. Keep tests parallelized and assert errors with `assert.Equal` / `assert.NotNil`.
- For UI logic (currently untested), follow the patterns established in other repos if/when tests are introduced; prefer React Testing Library + Vitest configured for Vite projects.

### Integration Testing

- Use `testcontainers-go` helpers (`tests/common_tests/container_builders`) to create disposable Postgres/Vault instances and clean them up in `AfterEach`-style helpers.
- Execute migrations through the same factories the services use (`orchestrator_db/factories`) before executing assertions, so schema drift is caught automatically.

## Technology-Specific Guidelines

### Go Services

- Every microservice is structured as a hexagonal/ports-and-adapters system: hosts and MVC controllers act as inbound adapters, `use_cases` hold the domain rules, `services` expose outbound ports, and `orchestrator_db` / HTTP clients / Vault / Redis packages are outbound adapters. Keep new code within the appropriate ring instead of short-circuiting layers.
- Bootstrap DI exactly as in `gateway/src/orchestrator_gateway_host/configurations/di_configuration.go`: register settings, then services/interfaces, then use cases. Use `dig.As` to bind interfaces instead of manual type assertions.
- HTTP servers are Fiber apps configured in `server_configuration.go` with `cors`, `fiberzap`, and Swagger routes. Add new controllers by implementing `ControllerInterface`, tagging them with `group:"controllers"`, and registering both authenticated/unauthenticated routes.
- Requests live in `use_cases/requests`, responses in `use_cases/responses`. Always add sanitize/validation helpers before hitting repositories.
- Repositories belong under `orchestrator_db/repositories` and must go through the shared `OrchestratorDbClientInterface` for transactions and migrations.
- For secrets/caching/blockchain:
  - Use `gateway_vault` services for KV storage paths (e.g., agent private keys follow `AgentPrivateKeyVaultPathTemplate`).
  - Use `gateway_redis/services/redis_storage.go` for cache interactions.
  - Use the helpers in `orchestrator_gateway/services/auth` for JWT/SIWE/signature flows; do not reimplement Keccak hashing.
- Migrations are timestamped Goose scripts in `*/migrations`. Update the Atlas schema file and rerun `generate_migrations.sh` when schema changes, and wire migrations via `ConfigureDatabase` so hosts auto-migrate at startup.
- Logging must go through the injected `*zap.Logger` (see `AuthenticationMandatoryMiddleware` and `agent_controller`). Never use `fmt.Println` outside of intentional CLI/debug utilities.

### TypeScript / React UI

- Respect the strict compiler options in `tsconfig.app.json` (no implicit `any`, `noUncheckedSideEffectImports`). New paths must be added under the `@/*` alias.
- Initialize Axios via `initAxiosDefaults` and always send JWT headers like the queries in `src/queries/api/agent.queries.ts`.
- Auth flows rely on `AuthProvider` (`contexts/auth-context.tsx`) storing the JWT in `localStorage`. Route guards redirect unauthenticated users back to `/login`. Hook into that context instead of duplicating login logic.
- Wallet connectivity flows through `@reown/appkit`: `App.tsx` picks the network based on `VITE_CHAIN_ID`, and `LoginForm` opens the modal via `useAppKit`. Copy that pattern for any new wallet interactions.
- Forms use `react-hook-form` + zod resolvers (see `components/dashboard/agents/forms/create-agent-first-step-form`). Keep schema files alongside the form and prefer asynchronous `superRefine` for on-chain validation.
- UI primitives live in `components/ui` and wrap Radix components. Extend them (e.g., add button variants via CVA) rather than hand-rolling new elements so focus states and theming remain coherent.
- Styling relies on Tailwind 4 utility classes merged through `cn` from `lib/utils.ts`; avoid inline styles except where SVG generation demands it.
- Toasts/errors route through `sonner` (`components/ui/sonner.tsx`).

### Docs (Docusaurus)

- Add or edit docs under `docs/docs/**`, keep front matter consistent, and register new pages in `docs/sidebars.ts`.
- Global config such as navbar links or Mermaid themes lives in `docs/docusaurus.config.ts`. Duplicate the existing patterns when exposing new sections.
- Run `npm run typecheck` inside `docs/` when you change TypeScript-based config/components.

### Image generation & Integrations

- The image generator (`agent-identity/images-generator`) composes SVGs in `services/agent_svg_image.go` and uses `chromedp` for screenshotting. Maintain deterministic string builders, and prefer the injected `AgentImageGenerationSettings` for tunables.
- Slack bridge and other connectors mirror the gateway layering: `*_host/main.go` → DI config → MVC → use cases → external HTTP clients. Follow their local modules when adding commands or event listeners.

## Version Control Guidelines

- No unified SemVer is enforced for the mono-repo; each deployable has its own release cadence. When you bump dependencies, update the relevant `go.sum`, `go.work.sum`, `package-lock.json` (UI/docs), or Docker tags inside the service directory.
- Goose migrations must remain chronological (`YYYYMMDDHHMMSS_description.sql`). Never edit a committed migration; add a new one instead.
- Regenerate mocks (`go generate` where configured) when interfaces change, and commit the updated files under `tests/**/mocks`.

## General Best Practices

- Run the relevant test suite(s) before publishing changes: `go test ./...` per service, `npm run lint && npm run build` inside `ui/`, and `npm run typecheck` inside `docs/`.
- Update `.env.example` files next to any service when new configuration is required.
- Keep Docker artifacts (`Dockerfile`, `docker-compose.yml`) synchronized with code changes so local environments continue to work.
- Prefer existing helper packages (`common_tests`, `services/interfaces`) over introducing new utility layers.

## Project-Specific Guidance

- The gateway orchestrates agent lifecycle operations: controllers (`orchestrator_gateway_mvc`) talk to use cases that coordinate blockchain (`orchestrator_blockchain`), database (`orchestrator_db`), Vault (`gateway_vault`), Redis (`gateway_redis`), and downstream HTTP connectors (`orchestrator_gateway_http_*`). Respect these seams; do not mix responsibilities.
- Agent Knowledge, Agent Reasoning, Agent Identity, and Slack Bridge services follow the same host/mvc/use-case/service/repository structure. When extending a capability, edit the matching service rather than the gateway.
- HTTP clients in `orchestrator_gateway_http_agent_*` already handle auth headers and DTOs for downstream agents. Extend those clients when new endpoints are exposed instead of inlining `http.Client` calls.
- Frontend routes in `ui/src/App.tsx` mirror backend capabilities (login, verify-wallet, dashboard sections). When adding a backend endpoint, expose it via the `queries/api` layer, wrap it in a context/hook if shared, and only then surface it in a page component.
- Documentation (`docs/docs/developers/**`, `docs/docs/users/**`) is structured by domain (Agent Identities, Reasoning, Knowledge, Channels). Keep new docs within the right domain folder so navigation stays predictable.
