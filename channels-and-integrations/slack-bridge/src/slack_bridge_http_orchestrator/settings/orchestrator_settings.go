package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const OrchestratorSettingsKey = "orchestrator"

type OrchestratorSettings struct {
	ApiBaseUrl string `mapstructure:"api_base_url"`
}

func NewOrchestratorSettings(logger *zap.Logger, v *viper.Viper) *OrchestratorSettings {
	if val, ok := os.LookupEnv("ORCHESTRATOR__API_BASE_URL"); ok {
		v.Set(OrchestratorSettingsKey+".api_base_url", val)
	}

	var settings OrchestratorSettings
	if err := v.UnmarshalKey(OrchestratorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal orchestrator settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.ApiBaseUrl)) == 0 {
		logger.Fatal("please set your orchestrator.api_base_url value in your configuration")
	}

	settings.ApiBaseUrl = strings.TrimSuffix(settings.ApiBaseUrl, "/")

	return &settings
}
