package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type VaultSettings struct {
	Token         string
	Address       string
	KvV1MountPath string
}

func NewVaultSettings(
	logger *zap.Logger,
) *VaultSettings {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if len(strings.TrimSpace(vaultAddr)) == 0 {
		logger.Fatal("Error, please set your VAULT_ADDR value in your environment")
	}

	token := os.Getenv("VAULT_TOKEN")
	if len(strings.TrimSpace(token)) == 0 {
		logger.Fatal("Error, please set your VAULT_TOKEN value in your environment")
	}

	kvV1MountPath := os.Getenv("VAULT_KV_V1_MOUNT_PATH")
	if len(strings.TrimSpace(kvV1MountPath)) == 0 {
		logger.Fatal("Error, please set your VAULT_KV_V1_MOUNT_PATH value in your environment")
	}

	return &VaultSettings{
		Token:         token,
		Address:       vaultAddr,
		KvV1MountPath: kvV1MountPath,
	}
}
