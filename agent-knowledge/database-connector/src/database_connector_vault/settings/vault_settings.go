package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const VaultSettingsKey = "vault"

type VaultSettings struct {
	Token         string `mapstructure:"token"`
	Address       string `mapstructure:"address"`
	KvV1MountPath string `mapstructure:"kv_v1_mount_path"`
}

func NewVaultSettings(
	logger *zap.Logger,
	v *viper.Viper,
) *VaultSettings {
	if val, ok := os.LookupEnv("VAULT__ADDRESS"); ok {
		v.Set(VaultSettingsKey+".address", val)
	}
	if val, ok := os.LookupEnv("VAULT__TOKEN"); ok {
		v.Set(VaultSettingsKey+".token", val)
	}
	if val, ok := os.LookupEnv("VAULT__KV_V1_MOUNT_PATH"); ok {
		v.Set(VaultSettingsKey+".kv_v1_mount_path", val)
	}

	var settings VaultSettings
	if err := v.UnmarshalKey(VaultSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal vault settings", zap.Error(err))
	}

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
