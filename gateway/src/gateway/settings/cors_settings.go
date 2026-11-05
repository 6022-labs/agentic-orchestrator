package settings

import (
	"net/url"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const CorsSettingsKey = "cors"

// CorsSettings lists the browser origins allowed to call the gateway; "*" allows any.
type CorsSettings struct {
	AllowOrigins []string `koanf:"allow_origins"`
}

func NewCorsSettings(logger *zap.Logger, k *koanf.Koanf, consoleSettings *ConsoleSettings) *CorsSettings {
	// Not seeded: mapstructure merges into an existing slice element-wise.
	var settings CorsSettings
	if err := k.Unmarshal(CorsSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal cors settings", zap.Error(err))
	}

	settings.normalizeOrigins(logger)
	if len(settings.AllowOrigins) == 0 {
		settings.AllowOrigins = []string{consoleSettings.Origin}
	}

	return &settings
}

// Env hands the whole list through as one comma-separated string.
func (s *CorsSettings) normalizeOrigins(logger *zap.Logger) {
	normalized := make([]string, 0, len(s.AllowOrigins))
	for _, entry := range s.AllowOrigins {
		for _, part := range strings.Split(entry, ",") {
			origin := strings.TrimRight(strings.TrimSpace(part), "/")
			if origin == "" {
				continue
			}
			if !s.isValidOrigin(origin) {
				logger.Fatal(
					"cors.allow_origins entries must be \"*\" or scheme://host[:port]",
					zap.String("origin", origin),
				)
			}
			normalized = append(normalized, origin)
		}
	}

	s.AllowOrigins = normalized
}

func (s *CorsSettings) isValidOrigin(origin string) bool {
	if origin == "*" {
		return true
	}
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Scheme != "" && parsed.Host != "" && parsed.Path == "" &&
		parsed.RawQuery == "" && parsed.Fragment == ""
}
