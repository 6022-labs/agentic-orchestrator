# Koanf Patterns (Generic)

## Bootstrap pattern

Use a single loader that merges file and env sources:

```go
func LoadKoanfConfig(configFile string) (*koanf.Koanf, error) {
    k := koanf.New(".")

    if configFile == "" {
        configFile = "appsettings.json"
    }

    if _, err := os.Stat(configFile); err == nil {
        if err := k.Load(file.Provider(configFile), json.Parser()); err != nil {
            return nil, err
        }
    } else if !errors.Is(err, os.ErrNotExist) {
        return nil, err
    }

    if err := loadEnvIntoKoanf(k); err != nil {
        return nil, err
    }

    return k, nil
}
```

## Environment key mapping

- Read only env vars containing `__`.
- Convert names to dotted lowercase keys.
- Examples:
  - `DATABASE__URL` -> `database.url`
  - `HOST__API_PORT` -> `host.api_port`

## Constructor pattern

```go
const DatabaseSettingsKey = "database"

type DatabaseSettings struct {
    URL string `koanf:"url"`
}

func NewDatabaseSettings(logger *zap.Logger, k *koanf.Koanf) *DatabaseSettings {
    settings := DatabaseSettings{}

    if err := k.Unmarshal(DatabaseSettingsKey, &settings); err != nil {
        logger.Fatal("failed to unmarshal database settings", zap.Error(err))
    }

    settings.URL = strings.TrimSpace(settings.URL)
    if settings.URL == "" {
        logger.Fatal("please set your database.url value in your configuration")
    }

    return &settings
}
```

## Validation rules

- Required field missing/invalid: `logger.Fatal`.
- Optional field missing/invalid: assign default and log `Warn`/`Info`.
- Normalize values before validation.
- Keep domain-specific validation explicit (URLs, durations, addresses, ranges).

## Test checklist

For each settings constructor, include at least:

1. **Load explicit values** from `confmap.Provider`.
2. **Default fallback behavior** for optional values.
3. **Environment override behavior** with `t.Setenv` for `__` keys.

Template for explicit values test:

```go
k := koanf.New(".")
_ = k.Load(confmap.Provider(map[string]interface{}{
    "database.url": "postgres://localhost:5432/app",
}, "."), nil)

settings := NewDatabaseSettings(zap.NewNop(), k)
assert.Equal(t, "postgres://localhost:5432/app", settings.URL)
```
