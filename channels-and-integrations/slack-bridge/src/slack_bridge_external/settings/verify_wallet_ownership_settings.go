package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const VerifyWalletOwnershipSettingsKey = "verify_wallet_ownership"

type VerifyWalletOwnershipSettings struct {
	OrchestratorFrontBaseUrl string `mapstructure:"orchestrator_front_base_url"`
}

func NewVerifyWalletOwnershipSettings(logger *zap.Logger, v *viper.Viper) *VerifyWalletOwnershipSettings {
	if val, ok := os.LookupEnv("VERIFY_WALLET_OWNERSHIP__ORCHESTRATOR_FRONT_BASE_URL"); ok {
		v.Set(VerifyWalletOwnershipSettingsKey+".orchestrator_front_base_url", val)
	}

	var settings VerifyWalletOwnershipSettings
	if err := v.UnmarshalKey(VerifyWalletOwnershipSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal verify wallet ownership settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.OrchestratorFrontBaseUrl)) == 0 {
		logger.Fatal("please set your verify_wallet_ownership.orchestrator_front_base_url value in your configuration")
	}

	settings.OrchestratorFrontBaseUrl = strings.TrimSuffix(settings.OrchestratorFrontBaseUrl, "/")

	return &settings
}
