package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const VaultSettingsKey = "vault"

type VaultSettings struct {
	Token         string `koanf:"token"`
	Address       string `koanf:"address"`
	KvV1MountPath string `koanf:"kv_v1_mount_path"`
}

func NewVaultSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *VaultSettings {
	settings := VaultSettings{}
	if err := k.Unmarshal(VaultSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal vault settings", zap.Error(err))
	}

	settings.Token = strings.TrimSpace(settings.Token)
	settings.Address = strings.TrimSpace(settings.Address)
	settings.KvV1MountPath = strings.TrimSpace(settings.KvV1MountPath)

	if len(strings.TrimSpace(settings.Address)) == 0 {
		logger.Fatal("please set your vault.address value in your configuration")
	}

	if len(strings.TrimSpace(settings.Token)) == 0 {
		logger.Fatal("please set your vault.token value in your configuration")
	}

	if len(strings.TrimSpace(settings.KvV1MountPath)) == 0 {
		logger.Fatal("please set your vault.kv_v1_mount_path value in your configuration")
	}

	return &settings
}
