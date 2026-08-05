package settings

import (
	"net/url"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SiweSettingsKey = "siwe"

// SiweSettings are the server-controlled fields of the EIP-4361 sign-in message.
type SiweSettings struct {
	Uri       string `koanf:"uri"`
	Statement string `koanf:"statement"`
	// Domain is the EIP-4361 authority shown by wallets, derived from Uri.
	Domain string `koanf:"-"`
}

func NewSiweSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *SiweSettings {
	settings := SiweSettings{
		Uri:       "http://localhost:5173",
		Statement: "Sign in to the 6022 orchestrator.",
	}
	if err := k.Unmarshal(SiweSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal siwe settings", zap.Error(err))
	}

	settings.Uri = strings.TrimSpace(settings.Uri)
	parsed, err := url.Parse(settings.Uri)
	if err != nil || parsed.Host == "" {
		logger.Fatal("siwe.uri must be an absolute URL (e.g. https://agent.example.com)")
	}
	settings.Domain = parsed.Host

	return &settings
}
