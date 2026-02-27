package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/settings"
)

type OrchestratorPrivateKeyFetcher struct {
	vaultSecretStorage   interfaces.VaultSecretStorageInterface
	orchestratorSettings *settings.OrchestratorSettings
}

func NewOrchestratorPrivateKeyFetcher(
	vaultSecretStorage interfaces.VaultSecretStorageInterface,
	orchestratorSettings *settings.OrchestratorSettings,
) *OrchestratorPrivateKeyFetcher {
	return &OrchestratorPrivateKeyFetcher{
		vaultSecretStorage:   vaultSecretStorage,
		orchestratorSettings: orchestratorSettings,
	}
}

func (f *OrchestratorPrivateKeyFetcher) FetchPrivateKey() (*string, error) {
	path := f.getOrchestratorPrivateKeyPath()

	privateKey, err := f.vaultSecretStorage.Read(path)
	if err != nil {
		return nil, err
	}

	return &privateKey, nil
}

func (f *OrchestratorPrivateKeyFetcher) getOrchestratorPrivateKeyPath() string {
	return fmt.Sprintf("orchestrator/%s/private_key", f.orchestratorSettings.PublicKey)
}
