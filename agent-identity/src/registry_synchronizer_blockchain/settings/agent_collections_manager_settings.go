package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

//go:embed agent_collections_managers.json
var agentCollectionsManagersJson []byte

type agentCollectionsManagerEntry struct {
	ChainId int64  `json:"chainId"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type agentCollectionsManagersFile struct {
	Managers []agentCollectionsManagerEntry `json:"managers"`
}

// AgentCollectionsManagerSettings holds the AgentCollectionsManager address for
// the configured chain, resolved at startup from the built-in registry by
// BlockchainSettings.ChainId.
type AgentCollectionsManagerSettings struct {
	SmartContractAddress common.Address
}

func NewAgentCollectionsManagerSettings(
	logger *zap.Logger,
	blockchainSettings *common_settings.BlockchainSettings,
) *AgentCollectionsManagerSettings {
	var f agentCollectionsManagersFile
	if err := json.Unmarshal(agentCollectionsManagersJson, &f); err != nil {
		logger.Fatal("failed to parse built-in agent_collections_managers.json", zap.Error(err))
	}

	for _, m := range f.Managers {
		if m.ChainId != blockchainSettings.ChainId {
			continue
		}

		addr := strings.TrimSpace(m.Address)
		if !common.IsHexAddress(addr) {
			logger.Fatal(
				"AgentCollectionsManager address is missing or invalid for the configured chain; fix agent_collections_managers.json and rebuild",
				zap.Int64("chainId", blockchainSettings.ChainId),
				zap.String("name", m.Name),
				zap.String("address", addr),
			)
		}

		logger.Info("resolved AgentCollectionsManager from built-in registry",
			zap.Int64("chainId", blockchainSettings.ChainId),
			zap.String("name", m.Name),
			zap.String("address", addr),
		)
		return &AgentCollectionsManagerSettings{SmartContractAddress: common.HexToAddress(addr)}
	}

	logger.Fatal(
		"no AgentCollectionsManager entry for the configured chain; add an entry to agent_collections_managers.json and rebuild",
		zap.Int64("chainId", blockchainSettings.ChainId),
	)
	return nil
}
