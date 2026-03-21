package services

import (
	"fmt"

	common_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
)

type OrchestratorPrivateKeyFetcher struct {
	vaultSecretStorage   common_interfaces.VaultSecretStorageInterface
	orchestratorSettings *settings.OrchestratorSettings
}

func NewOrchestratorPrivateKeyFetcher(
	vaultSecretStorage common_interfaces.VaultSecretStorageInterface,
	orchestratorSettings *settings.OrchestratorSettings,
) *OrchestratorPrivateKeyFetcher {
	return &OrchestratorPrivateKeyFetcher{
		vaultSecretStorage:   vaultSecretStorage,
		orchestratorSettings: orchestratorSettings,
	}
}

func (f *OrchestratorPrivateKeyFetcher) FetchPrivateKey() (*string, error) {
	privateKey, err := f.vaultSecretStorage.Read(f.getOrchestratorPrivateKeyPath())
	if err != nil {
		return nil, err
	}
	return &privateKey, nil
}

func (f *OrchestratorPrivateKeyFetcher) getOrchestratorPrivateKeyPath() string {
	return fmt.Sprintf("orchestrator/%s/private_key", f.orchestratorSettings.PublicKey)
}
