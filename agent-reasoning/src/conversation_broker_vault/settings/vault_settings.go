package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const VaultSettingsKey = "vault"

// VaultSettings configures the mandatory credential store backing the orchestrator wallet key.
type VaultSettings struct {
	Token         string `koanf:"token"`
	Address       string `koanf:"address"`
	KvV1MountPath string `koanf:"kv_v1_mount_path"`
}

func NewVaultSettings(logger *zap.Logger, k *koanf.Koanf) *VaultSettings {
	settings := VaultSettings{
		KvV1MountPath: "secret",
	}
	if err := k.Unmarshal(VaultSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal vault settings", zap.Error(err))
	}

	settings.Address = strings.TrimSpace(settings.Address)
	settings.Token = strings.TrimSpace(settings.Token)
	settings.KvV1MountPath = strings.TrimSpace(settings.KvV1MountPath)

	if len(settings.Address) == 0 {
		logger.Fatal("please set your vault.address value in your configuration")
	}
	if len(settings.Token) == 0 {
		logger.Fatal("please set your vault.token value in your configuration")
	}
	if len(settings.KvV1MountPath) == 0 {
		settings.KvV1MountPath = "secret"
	}

	return &settings
}
