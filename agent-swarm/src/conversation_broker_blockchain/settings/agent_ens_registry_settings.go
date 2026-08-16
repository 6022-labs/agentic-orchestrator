package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentEnsRegistrySettingsKey = "agent_ens_registry"

//go:embed agent_ens_registries.json
var agentEnsRegistriesJson []byte

type agentEnsRegistryEntry struct {
	ChainId int64  `json:"chainId"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type agentEnsRegistriesFile struct {
	Registries []agentEnsRegistryEntry `json:"registries"`
}

// AgentEnsRegistrySettings points at the contract L1 ENS resolves into; the broker reads it directly.
type AgentEnsRegistrySettings struct {
	SmartContractAddress common.Address `koanf:"smart_contract_address"`
}

func NewAgentEnsRegistrySettings(
	logger *zap.Logger,
	k *koanf.Koanf,
	blockchainSettings *BlockchainSettings,
) *AgentEnsRegistrySettings {
	var file agentEnsRegistriesFile
	if err := json.Unmarshal(agentEnsRegistriesJson, &file); err != nil {
		logger.Fatal("failed to parse built-in agent_ens_registries.json", zap.Error(err))
	}

	addressStr := strings.TrimSpace(k.String(AgentEnsRegistrySettingsKey + ".smart_contract_address"))
	if len(addressStr) == 0 {
		for _, registry := range file.Registries {
			if registry.ChainId == blockchainSettings.ChainId {
				addressStr = strings.TrimSpace(registry.Address)
				break
			}
		}
	}

	if len(addressStr) == 0 {
		logger.Fatal(
			"no built-in AgentEnsRegistry for this chain, "+
				"please set your agent_ens_registry.smart_contract_address value in your configuration",
			zap.Int64("chainId", blockchainSettings.ChainId),
		)
	}

	if !common.IsHexAddress(addressStr) {
		logger.Fatal(
			"agent_ens_registry.smart_contract_address is not a valid Ethereum address",
			zap.String("smartContractAddress", addressStr),
		)
	}

	return &AgentEnsRegistrySettings{SmartContractAddress: common.HexToAddress(addressStr)}
}
