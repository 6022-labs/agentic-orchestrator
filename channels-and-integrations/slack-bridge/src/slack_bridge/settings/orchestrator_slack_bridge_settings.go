package settings

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const OrchestratorSlackBridgeSettingsKey = "orchestrator_slack_bridge"

type OrchestratorSlackBridgeSettings struct {
	PublicKey common.Address `mapstructure:"public_key"`
}

func NewOrchestratorSlackBridgeSettings(logger *zap.Logger, v *viper.Viper) *OrchestratorSlackBridgeSettings {
	if val, ok := os.LookupEnv("ORCHESTRATOR_SLACK_BRIDGE__PUBLIC_KEY"); ok {
		v.Set(OrchestratorSlackBridgeSettingsKey+".public_key", val)
	}

	publicKeyStr := v.GetString(OrchestratorSlackBridgeSettingsKey + ".public_key")
	if len(strings.TrimSpace(publicKeyStr)) == 0 {
		logger.Fatal("please set your orchestrator_slack_bridge.public_key value in your configuration")
	}

	if !common.IsHexAddress(publicKeyStr) {
		logger.Fatal("orchestrator_slack_bridge.public_key is not a valid Ethereum address", zap.String("publicKey", publicKeyStr))
	}

	return &OrchestratorSlackBridgeSettings{
		PublicKey: common.HexToAddress(publicKeyStr),
	}
}
