package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

// BridgeCallbackSettings holds shared HTTP-client tuning for outbound calls
// to bridges. There is intentionally no `base_url` field — each bridge's
// base URL lives in the `bridges.base_url` column and is resolved per
// request from the database (the broker appends `/broker/callback`).
const BridgeCallbackSettingsKey = "bridge_callback"

type BridgeCallbackSettings struct {
	RequestTimeoutSeconds int `koanf:"request_timeout_seconds"`
}

func NewBridgeCallbackSettings(logger *zap.Logger, k *koanf.Koanf) *BridgeCallbackSettings {
	settings := BridgeCallbackSettings{
		RequestTimeoutSeconds: 30,
	}
	if err := k.Unmarshal(BridgeCallbackSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal bridge_callback settings", zap.Error(err))
	}

	if settings.RequestTimeoutSeconds <= 0 {
		logger.Warn("invalid bridge_callback.request_timeout_seconds, defaulting to 30")
		settings.RequestTimeoutSeconds = 30
	}

	return &settings
}
