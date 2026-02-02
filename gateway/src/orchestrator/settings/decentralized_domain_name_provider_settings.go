package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const DecentralizedDomainNameProviderSettingsKey = "decentralized_domain_name_provider"

type DecentralizedDomainNameProviderSettings struct {
	DomainExtension string `mapstructure:"domain_extension"`
}

func NewDecentralizedDomainNameProviderSettings(logger *zap.Logger, v *viper.Viper) *DecentralizedDomainNameProviderSettings {
	v.SetDefault(DecentralizedDomainNameProviderSettingsKey+".domain_extension", "eth")
	if val, ok := os.LookupEnv("DECENTRALIZED_DOMAIN__EXTENSION"); ok {
		v.Set(DecentralizedDomainNameProviderSettingsKey+".domain_extension", val)
	}

	var settings DecentralizedDomainNameProviderSettings
	if err := v.UnmarshalKey(DecentralizedDomainNameProviderSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal decentralized domain name provider settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.DomainExtension)) == 0 {
		logger.Info("no decentralized domain name provider domain extension set in configuration, using default 'eth'")
		settings.DomainExtension = v.GetString(DecentralizedDomainNameProviderSettingsKey + ".domain_extension")
	}

	return &settings
}
