package settings

import (
	"net/url"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ConsoleSettingsKey = "console"

// ConsoleSettings locate the web console; its URL signs users in and is the default CORS origin.
type ConsoleSettings struct {
	Url    string `koanf:"url"`
	Origin string `koanf:"-"`
	Domain string `koanf:"-"`
}

func NewConsoleSettings(logger *zap.Logger, k *koanf.Koanf) *ConsoleSettings {
	var settings ConsoleSettings
	if err := k.Unmarshal(ConsoleSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal console settings", zap.Error(err))
	}

	settings.Url = strings.TrimSpace(settings.Url)
	if settings.Url == "" {
		logger.Fatal("console.url is required")
	}
	parsed, err := url.Parse(settings.Url)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		logger.Fatal("console.url must be an absolute URL (e.g. https://orchestrator.example.com)")
	}
	settings.Origin = parsed.Scheme + "://" + parsed.Host
	settings.Domain = parsed.Host

	return &settings
}
