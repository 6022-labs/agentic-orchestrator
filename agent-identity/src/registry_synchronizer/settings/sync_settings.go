package settings

import (
	_ "embed"
	"encoding/json"
	"time"

	registry_common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SyncSettingsKey = "sync"

const (
	defaultBlockFinalizationDelay  = 2 * time.Second
	defaultBlockEventsPollInterval = 200 * time.Millisecond
)

// Public RPCs cap eth_getLogs at 10k blocks.
const defaultGapWindowBlocks = 10_000

//go:embed sync_start_blocks.json
var syncStartBlocksJson []byte

type syncStartBlockEntry struct {
	ChainId    int64  `json:"chainId"`
	Name       string `json:"name"`
	StartBlock uint64 `json:"startBlock"`
}

type syncStartBlocksFile struct {
	StartBlocks []syncStartBlockEntry `json:"startBlocks"`
}

type SyncSettings struct {
	StartBlock             uint64
	BlockFinalizationDelay time.Duration
	// BlockEventsPollInterval spaces the blocks taken from the queue, at most one each time.
	BlockEventsPollInterval time.Duration
	GapWindowBlocks         uint64
}

func NewSyncSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
	blockchainSettings *registry_common_settings.BlockchainSettings,
) *SyncSettings {
	var file syncStartBlocksFile
	if err := json.Unmarshal(syncStartBlocksJson, &file); err != nil {
		logger.Fatal("failed to parse built-in sync_start_blocks.json", zap.Error(err))
	}

	startBlockInt := k.Int64(SyncSettingsKey + ".start_block")
	if startBlockInt <= 0 {
		for _, entry := range file.StartBlocks {
			if entry.ChainId == blockchainSettings.ChainId {
				startBlockInt = int64(entry.StartBlock)
				break
			}
		}
	}

	if startBlockInt <= 0 {
		logger.Fatal(
			"no built-in start block for this chain, please set a valid sync.start_block value in your configuration",
			zap.Int64("chainId", blockchainSettings.ChainId),
		)
	}

	blockFinalizationDelay := k.Duration(SyncSettingsKey + ".block_finalization_delay")
	if blockFinalizationDelay <= 0 {
		blockFinalizationDelay = defaultBlockFinalizationDelay
	}

	blockEventsPollInterval := k.Duration(SyncSettingsKey + ".block_events_poll_interval")
	if blockEventsPollInterval <= 0 {
		blockEventsPollInterval = defaultBlockEventsPollInterval
	}

	gapWindowBlocks := uint64(defaultGapWindowBlocks)
	if k.Exists(SyncSettingsKey + ".gap_window_blocks") {
		windowBlocks := k.Int64(SyncSettingsKey + ".gap_window_blocks")
		if windowBlocks <= 0 {
			logger.Fatal("sync.gap_window_blocks must be a positive block count",
				zap.String("value", k.String(SyncSettingsKey+".gap_window_blocks")),
			)
		}
		gapWindowBlocks = uint64(windowBlocks)
	}

	return &SyncSettings{
		StartBlock:              uint64(startBlockInt),
		BlockFinalizationDelay:  blockFinalizationDelay,
		BlockEventsPollInterval: blockEventsPollInterval,
		GapWindowBlocks:         gapWindowBlocks,
	}
}
