package settings

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type Token6022Settings struct {
	SmartContractAddress common.Address
}

func NewToken6022Settings(logger *zap.Logger) *Token6022Settings {
	smartContractAddressStr := os.Getenv("TOKEN_6022_SMART_CONTRACT_ADDRESS")
	if !common.IsHexAddress(smartContractAddressStr) {
		logger.Fatal("please set your TOKEN_6022_SMART_CONTRACT_ADDRESS value in your environment")
	}

	smartContractAddress := common.HexToAddress(smartContractAddressStr)

	return &Token6022Settings{
		SmartContractAddress: smartContractAddress,
	}
}
