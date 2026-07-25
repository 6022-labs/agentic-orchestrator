package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const X402PaymentSettingsKey = "x402_payment"

// X402PaymentSettings tunes the free-call x402 signature lifetime, clock-drift tolerance and requirements cache.
type X402PaymentSettings struct {
	SignatureLifetimeSeconds   uint64 `koanf:"signature_lifetime_seconds"`
	ClockDriftToleranceSeconds uint64 `koanf:"clock_drift_tolerance_seconds"`
	ChallengeCacheTtlSeconds   uint64 `koanf:"challenge_cache_ttl_seconds"`
}

func NewX402PaymentSettings(logger *zap.Logger, k *koanf.Koanf) *X402PaymentSettings {
	settings := X402PaymentSettings{
		SignatureLifetimeSeconds:   300,
		ClockDriftToleranceSeconds: 60,
		ChallengeCacheTtlSeconds:   300,
	}
	if err := k.Unmarshal(X402PaymentSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal x402 payment settings", zap.Error(err))
	}

	if settings.SignatureLifetimeSeconds == 0 {
		settings.SignatureLifetimeSeconds = 300
	}
	if settings.ClockDriftToleranceSeconds == 0 {
		settings.ClockDriftToleranceSeconds = 60
	}
	if settings.ChallengeCacheTtlSeconds == 0 {
		settings.ChallengeCacheTtlSeconds = 300
	}

	return &settings
}
