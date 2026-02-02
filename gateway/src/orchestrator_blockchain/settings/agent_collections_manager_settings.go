package settings

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const AgentCollectionsManagerSettingsKey = "agent_collections_manager"

type AgentCollectionsManagerSettings struct {
	SmartContractAddress common.Address `mapstructure:"smart_contract_address"`
}

func NewAgentCollectionsManagerSettings(logger *zap.Logger, v *viper.Viper) *AgentCollectionsManagerSettings {
	if val, ok := os.LookupEnv("AGENT_COLLECTIONS_MANAGER__SMART_CONTRACT_ADDRESS"); ok {
		v.Set(AgentCollectionsManagerSettingsKey+".smart_contract_address", val)
	}

	addressStr := v.GetString(AgentCollectionsManagerSettingsKey + ".smart_contract_address")
	if len(strings.TrimSpace(addressStr)) == 0 {
		logger.Fatal("please set your agent_collections_manager.smart_contract_address value in your configuration")
	}

	if !common.IsHexAddress(addressStr) {
		logger.Fatal("agent_collections_manager.smart_contract_address is not a valid Ethereum address", zap.String("smartContractAddress", addressStr))
	}

	return &AgentCollectionsManagerSettings{
		SmartContractAddress: common.HexToAddress(addressStr),
	}
}
