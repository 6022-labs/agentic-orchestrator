---
name: koanf-settings
description: Create and update generic koanf-based Go settings modules with config file + environment layering, defaults, validation, DI wiring, and unit tests. Use when editing Go configuration code (for example */settings/*.go, */configurations/*.go) or adding tests for koanf unmarshalling and env overrides.
---

# Koanf Settings

Implement reusable koanf settings patterns for Go services.

## Workflow

1. Read `references/koanf-patterns.md`.
2. Implement or update a config bootstrapper:
   - initialize `koanf.New(".")`
   - load optional config file (`appsettings.json` or equivalent)
   - map env variables with `__` to dotted lowercase keys (`A__B__C` -> `a.b.c`)
3. Implement each settings constructor with this order:
   - set defaults in the struct value
   - `k.Unmarshal(<section>, &settings)`
   - normalize values (`strings.TrimSpace`, `strings.TrimSuffix`, numeric guards)
   - fail fast for missing required fields (`logger.Fatal`)
   - apply fallback defaults for optional fields (`logger.Warn` or `logger.Info`)
4. Wire constructors in DI (`container.Provide(settings.NewXSettings)`) and provide `*koanf.Koanf` once at composition root.
5. Add tests for each settings module:
   - loading explicit values from `confmap`
   - default fallback behavior
   - env override mapping (`SECTION__FIELD`)
6. Run targeted `go test` commands for touched settings packages.

## Conventions

- Keep keys section-based and snake_case (for example `database.url`, `host.api_port`).
- Keep constructor signatures consistent (`func NewXSettings(logger *zap.Logger, k *koanf.Koanf) *XSettings`).
- Keep fatal messages explicit (`please set your <section>.<field> value in your configuration`).
- Keep constructors focused on config parsing/validation only.

## Assets

Use `assets/generic-go-service/` as copy-ready examples:

- `settings_loader.go` and `settings_loader_test.go`
- `database_settings.go` and `database_settings_test.go`
- `appsettings.example.json`
