package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const AgentDatabaseConnectorSettingsKey = "agent_database_connector"

type AgentDatabaseConnectorSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewAgentDatabaseConnectorSettings(logger *zap.Logger, v *viper.Viper) *AgentDatabaseConnectorSettings {
	if val, ok := os.LookupEnv("AGENT_DATABASE_CONNECTOR__BASE_URL"); ok {
		v.Set(AgentDatabaseConnectorSettingsKey+".base_url", val)
	}

	var settings AgentDatabaseConnectorSettings
	if err := v.UnmarshalKey(AgentDatabaseConnectorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent database connector settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.BaseUrl)) == 0 {
		logger.Fatal("please set your agent_database_connector.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
