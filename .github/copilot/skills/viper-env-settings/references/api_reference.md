# Repository Settings Patterns

## Core pattern

```go
const HostSettingsKey = "host"

type HostSettings struct {
    ApiPort int    `mapstructure:"api_port"`
    AppName string `mapstructure:"app_name"`
}

func NewHostSettings(logger *zap.Logger, v *viper.Viper) *HostSettings {
    v.SetDefault(HostSettingsKey+".api_port", 3000)

    if val, ok := os.LookupEnv("HOST__API_PORT"); ok {
        v.Set(HostSettingsKey+".api_port", val)
    }
    if val, ok := os.LookupEnv("HOST__APP_NAME"); ok {
        v.Set(HostSettingsKey+".app_name", val)
    }

    var settings HostSettings
    if err := v.UnmarshalKey(HostSettingsKey, &settings); err != nil {
        logger.Fatal("failed to unmarshal host settings", zap.Error(err))
    }

    if len(strings.TrimSpace(settings.AppName)) == 0 {
        logger.Fatal("please set your host.app_name value in your configuration")
    }

    return &settings
}
```

## Naming rules

- `section.field_name` (viper/mapstructure) ↔ `SECTION__FIELD_NAME` (env).
- Keep section constants singular and stable (`host`, `database`, `vault`, `jwt`).
- Keep env override calls explicit per field; do not rely on implicit `AutomaticEnv`.

## Validation rules

- Fail fast on unmarshal errors with `logger.Fatal`.
- Fail fast for missing required strings with `strings.TrimSpace(...)`.
- Guard invalid numeric values and apply defaults exactly like nearby files.

## Test rules

- Place tests under the service test tree in a `settings` folder.
- Use subtests:
  - `Given settings provided through viper`
  - `Given settings provided through env`
- Use `t.Setenv(...)` for env override coverage.

## Good file references

- `gateway/src/orchestrator_gateway_host/settings/host_settings.go`
- `gateway/src/orchestrator_db/settings/database_settings.go`
- `gateway/src/orchestrator_sc_sync_host/settings/host_settings.go`
- `channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings/processing_settings.go`
