package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const BlockchainSettingsKey = "blockchain"

type BlockchainSettings struct {
	ChainId int64 `mapstructure:"chain_id"`
}

func NewBlockchainSettings(logger *zap.Logger, v *viper.Viper) *BlockchainSettings {
	if val, ok := os.LookupEnv("BLOCKCHAIN__CHAIN_ID"); ok {
		v.Set(BlockchainSettingsKey+".chain_id", val)
	}

	var settings BlockchainSettings
	if err := v.UnmarshalKey(BlockchainSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal blockchain settings", zap.Error(err))
	}

	if settings.ChainId <= 0 {
		logger.Fatal("please set your blockchain.chain_id value in your configuration")
	}

	return &settings
}
