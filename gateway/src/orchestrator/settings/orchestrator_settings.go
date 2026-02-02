package settings

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const OrchestratorSettingsKey = "orchestrator"

type OrchestratorSettings struct {
	PublicKey common.Address `mapstructure:"public_key"`
}

func NewOrchestratorSettings(logger *zap.Logger, v *viper.Viper) *OrchestratorSettings {
	if val, ok := os.LookupEnv("ORCHESTRATOR__PUBLIC_KEY"); ok {
		v.Set(OrchestratorSettingsKey+".public_key", val)
	}

	publicKeyStr := v.GetString(OrchestratorSettingsKey + ".public_key")
	if len(strings.TrimSpace(publicKeyStr)) == 0 {
		logger.Fatal("please set your orchestrator.public_key value in your configuration")
	}

	if !common.IsHexAddress(publicKeyStr) {
		logger.Fatal("orchestrator.public_key is not a valid Ethereum address", zap.String("publicKey", publicKeyStr))
	}

	return &OrchestratorSettings{
		PublicKey: common.HexToAddress(publicKeyStr),
	}
}
