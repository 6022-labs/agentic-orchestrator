package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const DecentralizedDomainNameProviderSettingsKey = "decentralized_domain_name_provider"

type DecentralizedDomainNameProviderSettings struct {
	DomainExtension string `koanf:"domain_extension"`
}

func NewDecentralizedDomainNameProviderSettings(logger *zap.Logger, k *koanf.Koanf) *DecentralizedDomainNameProviderSettings {
	settings := DecentralizedDomainNameProviderSettings{DomainExtension: "eth"}
	if err := k.Unmarshal(DecentralizedDomainNameProviderSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal decentralized domain name provider settings", zap.Error(err))
	}

	settings.DomainExtension = strings.TrimSpace(settings.DomainExtension)

	if len(settings.DomainExtension) == 0 {
		logger.Info("no decentralized domain name provider domain extension set in configuration, using default 'eth'")
		settings.DomainExtension = "eth"
	}

	return &settings
}
