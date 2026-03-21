package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const BlockchainSettingsKey = "blockchain"

type BlockchainSettings struct {
	ChainId int64 `koanf:"chain_id"`
}

func NewBlockchainSettings(logger *zap.Logger, k *koanf.Koanf) *BlockchainSettings {
	var settings BlockchainSettings
	if err := k.Unmarshal(BlockchainSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal blockchain settings", zap.Error(err))
	}

	if settings.ChainId <= 0 {
		logger.Fatal("please set your blockchain.chain_id value in your configuration")
	}
	return &settings
}
