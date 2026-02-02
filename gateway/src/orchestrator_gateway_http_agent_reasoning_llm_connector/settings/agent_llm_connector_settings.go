package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const AgentLLMConnectorSettingsKey = "agent_llm_connector"

type AgentLLMConnectorSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewAgentLLMConnectorSettings(logger *zap.Logger, v *viper.Viper) *AgentLLMConnectorSettings {
	if val, ok := os.LookupEnv("AGENT_LLM_CONNECTOR__BASE_URL"); ok {
		v.Set(AgentLLMConnectorSettingsKey+".base_url", val)
	}

	var settings AgentLLMConnectorSettings
	if err := v.UnmarshalKey(AgentLLMConnectorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent llm connector settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.BaseUrl)) == 0 {
		logger.Fatal("please set your agent_llm_connector.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
