package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

//go:embed token_6022.json
var token6022Json []byte

type token6022Entry struct {
	ChainId int64  `json:"chainId"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type token6022File struct {
	Tokens []token6022Entry `json:"tokens"`
}

// Token6022Settings holds the $6022 ERC-20 address for the configured chain,
// resolved at startup from the built-in registry by BlockchainSettings.ChainId.
type Token6022Settings struct {
	SmartContractAddress common.Address
}

func NewToken6022Settings(
	logger *zap.Logger,
	blockchainSettings *common_settings.BlockchainSettings,
) *Token6022Settings {
	var f token6022File
	if err := json.Unmarshal(token6022Json, &f); err != nil {
		logger.Fatal("failed to parse built-in token_6022.json", zap.Error(err))
	}

	for _, t := range f.Tokens {
		if t.ChainId != blockchainSettings.ChainId {
			continue
		}

		addr := strings.TrimSpace(t.Address)
		if !common.IsHexAddress(addr) {
			logger.Fatal(
				"$6022 token address is missing or invalid for the configured chain; fix token_6022.json and rebuild",
				zap.Int64("chainId", blockchainSettings.ChainId),
				zap.String("name", t.Name),
				zap.String("address", addr),
			)
		}

		logger.Info("resolved $6022 token from built-in registry",
			zap.Int64("chainId", blockchainSettings.ChainId),
			zap.String("name", t.Name),
			zap.String("address", addr),
		)
		return &Token6022Settings{SmartContractAddress: common.HexToAddress(addr)}
	}

	logger.Fatal(
		"no $6022 token entry for the configured chain; add an entry to token_6022.json and rebuild",
		zap.Int64("chainId", blockchainSettings.ChainId),
	)
	return nil
}
