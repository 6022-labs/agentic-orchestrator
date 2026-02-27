package settings

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentCollectionsManagerSettingsKey = "agent_collections_manager"

type AgentCollectionsManagerSettings struct {
	SmartContractAddress common.Address `koanf:"smart_contract_address"`
}

func NewAgentCollectionsManagerSettings(logger *zap.Logger, k *koanf.Koanf) *AgentCollectionsManagerSettings {
	addressStr := strings.TrimSpace(k.String(AgentCollectionsManagerSettingsKey + ".smart_contract_address"))
	if len(addressStr) == 0 {
		logger.Fatal("please set your agent_collections_manager.smart_contract_address value in your configuration")
	}

	if !common.IsHexAddress(addressStr) {
		logger.Fatal("agent_collections_manager.smart_contract_address is not a valid Ethereum address", zap.String("smartContractAddress", addressStr))
	}

	return &AgentCollectionsManagerSettings{
		SmartContractAddress: common.HexToAddress(addressStr),
	}
}
