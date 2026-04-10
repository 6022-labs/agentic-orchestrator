package services

import (
	"context"

	vault_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_vault/settings"
	vault "github.com/hashicorp/vault/api"
)

type VaultSecretStorage struct {
	vaultClient   *vault.Client
	vaultSettings *vault_settings.VaultSettings
}

func NewVaultSecretStorage(
	vaultClient *vault.Client,
	vaultSettings *vault_settings.VaultSettings,
) *VaultSecretStorage {
	return &VaultSecretStorage{
		vaultClient:   vaultClient,
		vaultSettings: vaultSettings,
	}
}

func (v *VaultSecretStorage) StoreKey(path, key, value string) error {
	return v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Put(context.Background(), path, map[string]interface{}{
			key: value,
		})
}

func (v *VaultSecretStorage) Store(path string, value map[string]string) error {
	data := make(map[string]interface{}, len(value))
	for k, v := range value {
		data[k] = v
	}
	return v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Put(context.Background(), path, data)
}

func (v *VaultSecretStorage) ReadKey(path, key string) (string, error) {
	secret, err := v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Get(context.Background(), path)
	if err != nil {
		return "", err
	}

	return secret.Data[key].(string), nil
}

func (v *VaultSecretStorage) Read(path string) (map[string]string, error) {
	secret, err := v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Get(context.Background(), path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		if strValue, ok := v.(string); ok {
			result[k] = strValue
		}
	}

	return result, nil
}

func (v *VaultSecretStorage) DeleteKey(path, key string) error {
	secret, err := v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Get(context.Background(), path)
	if err != nil {
		return err
	}

	delete(secret.Data, key)

	if len(secret.Data) == 0 {
		return v.vaultClient.
			KVv1(v.vaultSettings.KvV1MountPath).
			Delete(context.Background(), path)
	}

	return v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Put(context.Background(), path, secret.Data)
}

func (v *VaultSecretStorage) Delete(path string) error {
	return v.vaultClient.
		KVv1(v.vaultSettings.KvV1MountPath).
		Delete(context.Background(), path)
}
