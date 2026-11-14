package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_blockchain_integration_tests/fixtures"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// chainFixture is a throwaway chain plus a freshly deployed contract stack; the
// anvil account owns the minted agents, so it may publish their ENS texts.
type chainFixture struct {
	anvilAccountSettings container_builders.AnvilAccountSettings
	client               *ethclient.Client
	stack                *fixtures.AgentContractStack
}

func newChainFixture(t *testing.T) *chainFixture {
	anvilAccountSettings := container_builders.GetAnvilContainerBuilder().Create()

	client, err := ethclient.Dial(anvilAccountSettings.HttpUrl)
	if err != nil {
		t.Fatalf("failed to dial anvil: %v", err)
	}
	t.Cleanup(client.Close)

	stack, err := fixtures.DeployAgentStack(context.Background(), client, anvilAccountSettings)
	if err != nil {
		t.Fatalf("failed to deploy contract stack: %v", err)
	}

	return &chainFixture{
		anvilAccountSettings: anvilAccountSettings,
		client:               client,
		stack:                stack,
	}
}

func (f *chainFixture) chainId() string {
	return f.anvilAccountSettings.ChainId
}

// ensDomainFor is the name the broker resolves for an agent of this stack.
func (f *chainFixture) ensDomainFor(agentName string) string {
	return fixtures.EnsDomainFor(f.stack, f.chainId(), agentName)
}

// publishAgent mints an agent and publishes its url record; an empty url mints without provisioning.
func (f *chainFixture) publishAgent(t *testing.T, agentName, agentUrl string) {
	if err := fixtures.PublishAgent(
		context.Background(),
		f.client,
		f.anvilAccountSettings,
		f.stack,
		agentName,
		agentUrl,
	); err != nil {
		t.Fatalf("failed to publish %s: %v", agentName, err)
	}
}

func (f *chainFixture) ensUrlResolver() *services.EnsUrlResolver {
	return services.NewEnsUrlResolver(
		zap.NewNop(),
		f.client,
		&settings.BlockchainSettings{ChainId: int64(f.anvilAccountSettings.ChainIdUint)},
		&settings.RpcNodeSettings{RequestTimeoutSeconds: 30},
		&settings.AgentEnsRegistrySettings{SmartContractAddress: f.stack.EnsRegistry},
	)
}
