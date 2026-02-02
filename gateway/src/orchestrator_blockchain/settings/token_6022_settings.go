package settings

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const Token6022SettingsKey = "token_6022"

type Token6022Settings struct {
	SmartContractAddress common.Address `mapstructure:"smart_contract_address"`
}

func NewToken6022Settings(logger *zap.Logger, v *viper.Viper) *Token6022Settings {
	if val, ok := os.LookupEnv("TOKEN_6022__SMART_CONTRACT_ADDRESS"); ok {
		v.Set(Token6022SettingsKey+".smart_contract_address", val)
	}

	addressStr := v.GetString(Token6022SettingsKey + ".smart_contract_address")
	if len(strings.TrimSpace(addressStr)) == 0 {
		logger.Fatal("please set your token_6022.smart_contract_address value in your configuration")
	}

	if !common.IsHexAddress(addressStr) {
		logger.Fatal("token_6022.smart_contract_address is not a valid Ethereum address", zap.String("smartContractAddress", addressStr))
	}

	return &Token6022Settings{
		SmartContractAddress: common.HexToAddress(addressStr),
	}
}
