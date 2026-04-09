package settings

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const Token6022SettingsKey = "token_6022"

type Token6022Settings struct {
	SmartContractAddress common.Address `koanf:"smart_contract_address"`
}

func NewToken6022Settings(logger *zap.Logger, k *koanf.Koanf) *Token6022Settings {
	addressStr := strings.TrimSpace(k.String(Token6022SettingsKey + ".smart_contract_address"))
	if len(addressStr) == 0 {
		logger.Fatal("please set your token_6022.smart_contract_address value in your configuration")
	}

	if !common.IsHexAddress(addressStr) {
		logger.Fatal("token_6022.smart_contract_address is not a valid Ethereum address", zap.String("smartContractAddress", addressStr))
	}

	return &Token6022Settings{
		SmartContractAddress: common.HexToAddress(addressStr),
	}
}
