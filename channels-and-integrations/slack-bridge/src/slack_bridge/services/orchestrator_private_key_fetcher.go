package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
)

type OrchestratorPrivateKeyFetcher struct {
	vaultSecretStorage              interfaces.VaultSecretStorageInterface
	orchestratorSlackBridgeSettings *settings.OrchestratorSlackBridgeSettings
}

func NewOrchestratorPrivateKeyFetcher(
	vaultSecretStorage interfaces.VaultSecretStorageInterface,
	orchestratorSlackBridgeSettings *settings.OrchestratorSlackBridgeSettings,
) *OrchestratorPrivateKeyFetcher {
	return &OrchestratorPrivateKeyFetcher{
		vaultSecretStorage:              vaultSecretStorage,
		orchestratorSlackBridgeSettings: orchestratorSlackBridgeSettings,
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
	return fmt.Sprintf("orchestrator/%s/private_key", f.orchestratorSlackBridgeSettings.PublicKey)
}
