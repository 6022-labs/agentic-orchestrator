package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	registry_common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentCollectionsManagerSettingsKey = "agent_collections_manager"

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

type AgentCollectionsManagerSettings struct {
	SmartContractAddress common.Address `koanf:"smart_contract_address"`
}

func NewAgentCollectionsManagerSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
	blockchainSettings *registry_common_settings.BlockchainSettings,
) *AgentCollectionsManagerSettings {
	var file agentCollectionsManagersFile
	if err := json.Unmarshal(agentCollectionsManagersJson, &file); err != nil {
		logger.Fatal("failed to parse built-in agent_collections_managers.json", zap.Error(err))
	}

	addressStr := strings.TrimSpace(k.String(AgentCollectionsManagerSettingsKey + ".smart_contract_address"))
	if len(addressStr) == 0 {
		for _, manager := range file.Managers {
			if manager.ChainId == blockchainSettings.ChainId {
				addressStr = strings.TrimSpace(manager.Address)
				break
			}
		}
	}

	if len(addressStr) == 0 {
		logger.Fatal("no built-in AgentCollectionsManager for this chain, please set your agent_collections_manager.smart_contract_address value in your configuration", zap.Int64("chainId", blockchainSettings.ChainId))
	}

	if !common.IsHexAddress(addressStr) {
		logger.Fatal("agent_collections_manager.smart_contract_address is not a valid Ethereum address", zap.String("smartContractAddress", addressStr))
	}

	return &AgentCollectionsManagerSettings{
		SmartContractAddress: common.HexToAddress(addressStr),
	}
}
