package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SyncSettingsKey = "sync"

type SyncSettings struct {
	StartBlock uint64
}

func NewSyncSettings(logger *zap.Logger, k *koanf.Koanf) *SyncSettings {
	startBlock := uint64(k.Int64(SyncSettingsKey + ".start_block"))

	if startBlock == 0 {
		logger.Fatal("please set a valid sync.start_block value in your configuration")
	}

	return &SyncSettings{
		StartBlock: startBlock,
	}
}
