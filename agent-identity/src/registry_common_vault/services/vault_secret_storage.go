package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_vault/settings"
	vault "github.com/hashicorp/vault/api"
)

type VaultSecretStorage struct {
	vaultClient   *vault.Client
	vaultSettings *settings.VaultSettings
}

func NewVaultSecretStorage(
	vaultClient *vault.Client,
	vaultSettings *settings.VaultSettings,
) *VaultSecretStorage {
	return &VaultSecretStorage{
		vaultClient:   vaultClient,
		vaultSettings: vaultSettings,
	}
}

func (v *VaultSecretStorage) Store(key string, value string) error {
	return v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Put(context.Background(), key, map[string]interface{}{
			"value": value,
		})
}

func (v *VaultSecretStorage) Read(key string) (string, error) {
	secret, err := v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	value, ok := secret.Data["value"].(string)
	if !ok {
		return "", fmt.Errorf("secret key %q: missing or non-string \"value\" field", key)
	}

	return value, nil
}

func (v *VaultSecretStorage) Delete(key string) error {
	err := v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Delete(context.Background(), key)
	if err != nil {
		return err
	}

	return nil
}
