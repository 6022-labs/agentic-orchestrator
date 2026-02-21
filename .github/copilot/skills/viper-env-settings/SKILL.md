---
name: viper-env-settings
description: Create or update Go settings modules in this repo so they work with both Viper config keys and environment variable overrides. Use when editing files under */settings/*.go, adding new service configuration, wiring defaults/validation, or writing settings tests that must pass for both viper-provided values and env values.
---

# Viper Env Settings

Implement repository-consistent settings loaders that read from Viper, support explicit env overrides, and validate required values.

Use this workflow every time you add or modify settings in Go services.

## Workflow

1. Pick a nearby settings file in the same service and mirror its structure.
   - Examples: `gateway/src/orchestrator_db/settings/database_settings.go`, `gateway/src/orchestrator_gateway_host/settings/host_settings.go`.
2. Define a `...SettingsKey` constant and a struct with `mapstructure` tags.
3. Set defaults for optional fields with `v.SetDefault`.
4. Add env overrides with `os.LookupEnv` + `v.Set` for each field.
   - Convert `section.field_name` to `SECTION__FIELD_NAME`.
5. Unmarshal with `v.UnmarshalKey(..., &settings)` and fail fast with `logger.Fatal` on unmarshal errors.
6. Validate required values (`strings.TrimSpace`, numeric checks, or domain checks) and fail fast for invalid required config.
7. Add/adjust unit tests in the matching `tests/**/settings/when_loading_*_settings_test.go` file:
   - Case 1: values provided through `v.Set(...)`.
   - Case 2: env overrides with `t.Setenv(...)`.
8. Run only the relevant settings tests for the touched module.

## Conventions To Enforce

- Keep env variable names uppercase with `__` separators (example: `HOST__API_PORT`).
- Keep keys and tags aligned:
  - `v.Set("host.api_port", value)`
  - `` `mapstructure:"api_port"` ``
- Keep constructor signature `func NewXSettings(logger *zap.Logger, v *viper.Viper) *XSettings`.
- Keep fatal messages explicit: `please set your <section>.<field> value in your configuration`.

## Resources

- Read `references/api_reference.md` for concrete repository patterns and snippets.
- Run `scripts/env_key_helper.py "<section.field_name>"` when you need a quick env variable name conversion.
