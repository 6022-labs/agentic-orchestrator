package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_blockchain_integration_tests/fixtures"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// chainFixture is a throwaway anvil chain with a deployed 6022 stack. The http
// and ws clients mirror the two named clients the host wires.
type chainFixture struct {
	t                    *testing.T
	anvilAccountSettings container_builders.AnvilAccountSettings
	stack                *fixtures.AgentContractStack
	httpClient           *ethclient.Client
	wsClient             *ethclient.Client
}

func newChainFixture(t *testing.T) *chainFixture {
	anvilAccountSettings := container_builders.CreateAnvilAccount()

	httpClient, err := ethclient.Dial(anvilAccountSettings.HttpUrl)
	if err != nil {
		t.Fatalf("failed to dial anvil over http: %v", err)
	}
	t.Cleanup(httpClient.Close)

	wsClient, err := ethclient.Dial(anvilAccountSettings.WsUrl)
	if err != nil {
		t.Fatalf("failed to dial anvil over websocket: %v", err)
	}
	t.Cleanup(wsClient.Close)

	stack, err := fixtures.DeployAgentStack(context.Background(), httpClient, anvilAccountSettings)
	if err != nil {
		t.Fatalf("failed to deploy contract stack: %v", err)
	}

	return &chainFixture{
		t:                    t,
		anvilAccountSettings: anvilAccountSettings,
		stack:                stack,
		httpClient:           httpClient,
		wsClient:             wsClient,
	}
}

func (f *chainFixture) chainId() string {
	return f.anvilAccountSettings.ChainId
}

func (f *chainFixture) agentCollectionsManagerSettings() *settings.AgentCollectionsManagerSettings {
	return &settings.AgentCollectionsManagerSettings{SmartContractAddress: f.stack.Manager}
}

func (f *chainFixture) blockchainRequester() *services.BlockchainRequester {
	return services.NewBlockchainRequester(services.NewBlockchainRequesterParams{
		Client: f.httpClient,
	})
}

func (f *chainFixture) agentCollectionRequester() *services.AgentCollectionRequester {
	return services.NewAgentCollectionRequester(services.NewAgentCollectionRequesterParams{
		Client: f.httpClient,
	})
}

func (f *chainFixture) agentCollectionsManagerRequester() *services.AgentCollectionsManagerRequester {
	return services.NewAgentCollectionsManagerRequester(services.NewAgentCollectionsManagerRequesterParams{
		Client:                          f.httpClient,
		AgentCollectionsManagerSettings: f.agentCollectionsManagerSettings(),
	})
}

func (f *chainFixture) agentCollectionEventSubscriptionProvider() *services.AgentCollectionEventSubscriptionProvider {
	return services.NewAgentCollectionEventSubscriptionProvider(services.NewAgentCollectionEventSubscriptionProviderParams{
		Logger: zap.NewNop(),
		Client: f.wsClient,
	})
}

func (f *chainFixture) agentCollectionsManagerEventSubscriptionProvider() *services.AgentCollectionsManagerEventSubscriptionProvider {
	return services.NewAgentCollectionsManagerEventSubscriptionProvider(services.NewAgentCollectionsManagerEventSubscriptionProviderParams{
		Logger:                          zap.NewNop(),
		Client:                          f.wsClient,
		AgentCollectionsManagerSettings: f.agentCollectionsManagerSettings(),
	})
}

func (f *chainFixture) currentBlock() uint64 {
	blockNumber, err := f.httpClient.BlockNumber(context.Background())
	if err != nil {
		f.t.Fatalf("failed to read the head block: %v", err)
	}
	return blockNumber
}

// mintAgent gives every agent its own evm address; the collection rejects a
// second agent claiming an address it already knows.
func (f *chainFixture) mintAgent(agentName string) {
	if err := fixtures.MintAgent(
		context.Background(),
		f.httpClient,
		f.anvilAccountSettings,
		f.stack.Collection,
		agentName,
		common_tests.RandomEthAddress(),
	); err != nil {
		f.t.Fatalf("failed to mint %s: %v", agentName, err)
	}
}
