package settings

import (
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SyncSettingsKey = "sync"

type SyncSettings struct {
	StartBlock             uint64
	BlockFinalizationDelay time.Duration
}

func NewSyncSettings(logger *zap.Logger, k *koanf.Koanf) *SyncSettings {
	startBlockInt := k.Int64(SyncSettingsKey + ".start_block")

	if startBlockInt <= 0 {
		logger.Fatal("please set a valid sync.start_block value in your configuration")
	}

	startBlock := uint64(startBlockInt)

	blockFinalizationDelay := k.Duration(SyncSettingsKey + ".block_finalization_delay")

	return &SyncSettings{
		StartBlock:             startBlock,
		BlockFinalizationDelay: blockFinalizationDelay,
	}
}
