package settings

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type OrchestratorSettings struct {
	PublicKey common.Address
}

func NewOrchestratorSettings(logger *zap.Logger) *OrchestratorSettings {
	publicKeyStr := os.Getenv("PUBLIC_KEY")
	if !common.IsHexAddress(publicKeyStr) {
		logger.Fatal("PUBLIC_KEY is not a valid Ethereum address", zap.String("publicKey", publicKeyStr))
	}

	publicKey := common.HexToAddress(publicKeyStr)

	return &OrchestratorSettings{
		PublicKey: publicKey,
	}
}
