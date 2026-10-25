package settings

import (
	"net/url"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SiweSettingsKey = "siwe"

const defaultSignedRequestMaxAgeSeconds = 60

// SiweSettings are the server-controlled fields of the EIP-4361 sign-in message.
type SiweSettings struct {
	Uri                        string `koanf:"uri"`
	Statement                  string `koanf:"statement"`
	SignedRequestMaxAgeSeconds int    `koanf:"signed_request_max_age_seconds"`
	// Domain is the EIP-4361 authority shown by wallets, derived from Uri.
	Domain string `koanf:"-"`
}

func NewSiweSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *SiweSettings {
	settings := SiweSettings{
		Uri:                        "http://localhost:5173",
		Statement:                  "Sign in to the 6022 orchestrator.",
		SignedRequestMaxAgeSeconds: defaultSignedRequestMaxAgeSeconds,
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

	if settings.SignedRequestMaxAgeSeconds <= 0 {
		logger.Warn("invalid siwe.signed_request_max_age_seconds value, defaulting to 60")
		settings.SignedRequestMaxAgeSeconds = defaultSignedRequestMaxAgeSeconds
	}

	return &settings
}
