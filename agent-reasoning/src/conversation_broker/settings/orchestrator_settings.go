package settings

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const OrchestratorSettingsKey = "orchestrator"

// OrchestratorSettings identifies the orchestrator wallet; the public key is mandatory.
type OrchestratorSettings struct {
	PublicKey common.Address `koanf:"public_key"`
}

func NewOrchestratorSettings(logger *zap.Logger, k *koanf.Koanf) *OrchestratorSettings {
	publicKeyStr := strings.TrimSpace(k.String(OrchestratorSettingsKey + ".public_key"))
	if len(publicKeyStr) == 0 {
		logger.Fatal("please set your orchestrator.public_key value in your configuration")
	}
	if !common.IsHexAddress(publicKeyStr) {
		logger.Fatal("orchestrator.public_key is not a valid Ethereum address", zap.String("publicKey", publicKeyStr))
	}

	settings := OrchestratorSettings{}
	if err := k.Unmarshal(OrchestratorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal orchestrator settings", zap.Error(err))
	}

	return &settings
}
