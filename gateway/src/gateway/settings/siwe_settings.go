package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SiweSettingsKey = "siwe"

const defaultSignedRequestMaxAgeSeconds = 60

// SiweSettings are the server-controlled fields of the EIP-4361 sign-in message.
type SiweSettings struct {
	Statement                  string `koanf:"statement"`
	SignedRequestMaxAgeSeconds int    `koanf:"signed_request_max_age_seconds"`
	// Uri and Domain come from console.url; Domain is the authority wallets show.
	Uri    string `koanf:"-"`
	Domain string `koanf:"-"`
}

func NewSiweSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
	consoleSettings *ConsoleSettings,
) *SiweSettings {
	settings := SiweSettings{
		Statement:                  "Sign in to the 6022 orchestrator.",
		SignedRequestMaxAgeSeconds: defaultSignedRequestMaxAgeSeconds,
	}
	if err := k.Unmarshal(SiweSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal siwe settings", zap.Error(err))
	}

	settings.Uri = consoleSettings.Url
	settings.Domain = consoleSettings.Domain

	if settings.SignedRequestMaxAgeSeconds <= 0 {
		logger.Warn("invalid siwe.signed_request_max_age_seconds value, defaulting to 60")
		settings.SignedRequestMaxAgeSeconds = defaultSignedRequestMaxAgeSeconds
	}

	return &settings
}
