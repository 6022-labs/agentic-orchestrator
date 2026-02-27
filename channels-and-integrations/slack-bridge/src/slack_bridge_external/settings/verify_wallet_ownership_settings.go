package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const VerifyWalletOwnershipSettingsKey = "verify_wallet_ownership"

type VerifyWalletOwnershipSettings struct {
	OrchestratorFrontBaseUrl string `koanf:"orchestrator_front_base_url"`
}

func NewVerifyWalletOwnershipSettings(logger *zap.Logger, k *koanf.Koanf) *VerifyWalletOwnershipSettings {
	settings := VerifyWalletOwnershipSettings{}
	if err := k.Unmarshal(VerifyWalletOwnershipSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal verify wallet ownership settings", zap.Error(err))
	}

	settings.OrchestratorFrontBaseUrl = strings.TrimSpace(settings.OrchestratorFrontBaseUrl)

	if len(strings.TrimSpace(settings.OrchestratorFrontBaseUrl)) == 0 {
		logger.Fatal("please set your verify_wallet_ownership.orchestrator_front_base_url value in your configuration")
	}

	settings.OrchestratorFrontBaseUrl = strings.TrimSuffix(settings.OrchestratorFrontBaseUrl, "/")

	return &settings
}
