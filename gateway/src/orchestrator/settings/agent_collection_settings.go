package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const AgentCollectionSettingsKey = "agent_collection"

type AgentCollectionSettings struct {
	SigningDomain string `mapstructure:"signing_domain"`
}

func NewAgentCollectionSettings(logger *zap.Logger, v *viper.Viper) *AgentCollectionSettings {
	if val, ok := os.LookupEnv("AGENT_COLLECTION__SIGNING_DOMAIN"); ok {
		v.Set(AgentCollectionSettingsKey+".signing_domain", val)
	}

	var settings AgentCollectionSettings
	if err := v.UnmarshalKey(AgentCollectionSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent collection settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.SigningDomain)) == 0 {
		logger.Fatal("please set your agent_collection.signing_domain value in your configuration")
	}

	return &settings
}
