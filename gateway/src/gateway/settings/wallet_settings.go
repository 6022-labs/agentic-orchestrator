package settings

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const WalletSettingsKey = "wallet"

// WalletSettings holds the orchestrator signing key; the address is derived once at load.
type WalletSettings struct {
	PrivateKey string         `koanf:"private_key"`
	Address    common.Address `koanf:"-"`
}

func NewWalletSettings(logger *zap.Logger, k *koanf.Koanf) *WalletSettings {
	settings := WalletSettings{}
	if err := k.Unmarshal(WalletSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal wallet settings", zap.Error(err))
	}

	settings.PrivateKey = strings.TrimSpace(settings.PrivateKey)
	if len(settings.PrivateKey) == 0 {
		logger.Fatal("please set your wallet.private_key value in your configuration")
	}

	key, err := crypto.HexToECDSA(strings.TrimPrefix(settings.PrivateKey, "0x"))
	if err != nil {
		logger.Fatal("wallet.private_key is not a valid private key", zap.Error(err))
	}
	settings.Address = crypto.PubkeyToAddress(key.PublicKey)

	return &settings
}
