package settings

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type BlockchainSettings struct {
	ChainId int64
}

func NewBlockchainSettings(logger *zap.Logger) *BlockchainSettings {
	chainIdStr := os.Getenv("CHAIN_ID")
	if len(strings.TrimSpace(chainIdStr)) == 0 {
		logger.Fatal("please set your CHAIN_ID value in your environment")
	}

	chainId, err := strconv.Atoi(chainIdStr)
	if err != nil {
		logger.Fatal("CHAIN_ID must be an integer")
	}

	if chainId < 0 {
		logger.Fatal("CHAIN_ID must be a positive integer")
	}

	return &BlockchainSettings{
		ChainId: int64(chainId),
	}
}
