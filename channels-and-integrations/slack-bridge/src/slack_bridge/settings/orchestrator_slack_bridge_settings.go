package settings

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type OrchestratorSlackBridgeSettings struct {
	PublicKey common.Address
}

func NewOrchestratorSlackBridgeSettings(logger *zap.Logger) *OrchestratorSlackBridgeSettings {
	publicKeyStr := os.Getenv("PUBLIC_KEY")

	if !common.IsHexAddress(publicKeyStr) {
		logger.Fatal("PUBLIC_KEY is not a valid Ethereum address", zap.String("publicKey", publicKeyStr))
	}

	publicKey := common.HexToAddress(publicKeyStr)

	return &OrchestratorSlackBridgeSettings{
		PublicKey: publicKey,
	}
}
