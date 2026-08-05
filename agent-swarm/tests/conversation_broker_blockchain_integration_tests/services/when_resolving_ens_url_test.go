package services_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_blockchain_integration_tests/fixtures"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	resolvedEnsDomain = "agent.6022.eth"
	resolvedAgentUrl  = "https://agent.example.com"
)

type WhenResolvingEnsUrlTestingSuite struct {
	anvilAccountSettings container_builders.AnvilAccountSettings
	client               *ethclient.Client
	gateway              *fixtures.EnsGatewayStub
	stack                *fixtures.EnsContractStack
	sut                  *services.EnsUrlResolver
}

// WhenResolvingEnsUrlBeforeEach claims registeredDomain in the registry and points
// it at a real offchain
// resolver whose proofs the gateway stub signs.
func WhenResolvingEnsUrlBeforeEach(t *testing.T, registeredDomain string) *WhenResolvingEnsUrlTestingSuite {
	anvilAccountSettings := container_builders.CreateAnvilAccount()

	client, err := ethclient.Dial(anvilAccountSettings.HttpUrl)
	if err != nil {
		t.Fatalf("failed to dial anvil: %v", err)
	}
	t.Cleanup(client.Close)

	encodedResult, err := fixtures.EncodeStringResult(resolvedAgentUrl)
	if err != nil {
		t.Fatalf("failed to encode the gateway result: %v", err)
	}
	gateway, err := fixtures.NewEnsGatewayStub(encodedResult)
	if err != nil {
		t.Fatalf("failed to start the gateway stub: %v", err)
	}
	t.Cleanup(gateway.Close)

	stack, err := fixtures.DeployEnsStack(
		context.Background(),
		client,
		anvilAccountSettings,
		[]string{gateway.Server.URL},
		gateway.SignerAddr,
		registeredDomain,
	)
	if err != nil {
		t.Fatalf("failed to deploy the ENS stack: %v", err)
	}
	// The proof is bound to the resolver address, so it is only known post-deploy.
	gateway.Target = stack.Resolver

	return &WhenResolvingEnsUrlTestingSuite{
		anvilAccountSettings: anvilAccountSettings,
		client:               client,
		gateway:              gateway,
		stack:                stack,
		sut: services.NewEnsUrlResolver(services.NewEnsUrlResolverParams{
			Logger:     zap.NewNop(),
			EthCaller:  client,
			HttpClient: &http.Client{Timeout: 10 * time.Second},
			EnsResolverSettings: &settings.EnsResolverSettings{
				RpcUrl:                       anvilAccountSettings.HttpUrl,
				RegistryAddress:              stack.Registry.Hex(),
				GatewayRequestTimeoutSeconds: 10,
			},
		}),
	}
}

func TestWhenResolvingEnsUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given the name itself is registered with an offchain resolver", func(t *testing.T) {
		t.Parallel()

		t.Run("Should follow the OffchainLookup revert and return the proven url", func(t *testing.T) {
			t.Parallel()
			suite := WhenResolvingEnsUrlBeforeEach(t, resolvedEnsDomain)

			url, err := suite.sut.ResolveUrl(resolvedEnsDomain)

			assert.NoError(t, err)
			assert.Equal(t, resolvedAgentUrl, url)
		})
	})

	t.Run("Given only the parent domain carries the resolver", func(t *testing.T) {
		t.Parallel()

		t.Run("Should walk up the labels until it finds one", func(t *testing.T) {
			t.Parallel()
			suite := WhenResolvingEnsUrlBeforeEach(t, "6022.eth")

			url, err := suite.sut.ResolveUrl(resolvedEnsDomain)

			assert.NoError(t, err)
			assert.Equal(t, resolvedAgentUrl, url)
		})
	})

	t.Run("Given no resolver is registered for the name or any parent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail instead of calling the zero address", func(t *testing.T) {
			t.Parallel()
			suite := WhenResolvingEnsUrlBeforeEach(t, "someone-else.eth")

			_, err := suite.sut.ResolveUrl(resolvedEnsDomain)

			assert.Error(t, err)
		})
	})

	t.Run("Given the gateway signs a proof the resolver does not trust", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail on-chain verification rather than trust the answer", func(t *testing.T) {
			t.Parallel()
			suite := WhenResolvingEnsUrlBeforeEach(t, resolvedEnsDomain)
			suite.gateway.CorruptSign = true

			_, err := suite.sut.ResolveUrl(resolvedEnsDomain)

			assert.Error(t, err)
		})
	})

	t.Run("Given the registry points the name at a contract that is not a resolver", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail instead of treating the response as a lookup", func(t *testing.T) {
			t.Parallel()
			suite := WhenResolvingEnsUrlBeforeEach(t, resolvedEnsDomain)

			// The registry contract answers, but not with an OffchainLookup revert.
			err := fixtures.SetResolver(
				context.Background(),
				suite.client,
				suite.anvilAccountSettings,
				suite.stack.Registry,
				resolvedEnsDomain,
				suite.stack.Registry,
			)
			assert.NoError(t, err)

			_, err = suite.sut.ResolveUrl(resolvedEnsDomain)

			assert.Error(t, err)
		})
	})

	t.Run("Given a resolver address with no contract code", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail rather than report an empty url", func(t *testing.T) {
			t.Parallel()
			suite := WhenResolvingEnsUrlBeforeEach(t, resolvedEnsDomain)

			err := fixtures.SetResolver(
				context.Background(),
				suite.client,
				suite.anvilAccountSettings,
				suite.stack.Registry,
				resolvedEnsDomain,
				common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
			)
			assert.NoError(t, err)

			_, err = suite.sut.ResolveUrl(resolvedEnsDomain)

			assert.Error(t, err)
		})
	})
}
