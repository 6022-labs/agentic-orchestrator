package services_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	contract_abi "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/helpers"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_blockchain_mocks/clients_mocks"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	agentEnsRegistryAddressForTest = "0xa3f6011779B55BB155BDAcd25Dd5Bcb9f5597646"
	resolvableEnsDomain            = "agent.collection.80002.6022.eth"
)

type WhenResolvingEnsUrlTestingSuite struct {
	sut *services.EnsUrlResolver

	mockEthClient *clients_mocks.MockEthClientInterface
	sentCall      ethereum.CallMsg
}

func WhenResolvingEnsUrlBeforeEach(t *testing.T) *WhenResolvingEnsUrlTestingSuite {
	mockController := gomock.NewController(t)

	mockEthClient := clients_mocks.NewMockEthClientInterface(mockController)

	return &WhenResolvingEnsUrlTestingSuite{
		sut: services.NewEnsUrlResolver(
			zap.NewNop(),
			mockEthClient,
			&settings.BlockchainSettings{ChainId: 80002},
			&settings.RpcNodeSettings{RequestTimeoutSeconds: 10},
			&settings.AgentEnsRegistrySettings{
				SmartContractAddress: common.HexToAddress(agentEnsRegistryAddressForTest),
			},
		),

		mockEthClient: mockEthClient,
	}
}

func packedTextReturn(t *testing.T, url string) []byte {
	stringType, err := abi.NewType("string", "", nil)
	assert.NoError(t, err)

	returnData, err := abi.Arguments{{Type: stringType}}.Pack(url)
	assert.NoError(t, err)

	return returnData
}

func TestWhenResolvingEnsUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given the registry holds a url text record", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingEnsUrlTestingSuite) {
			suite.mockEthClient.EXPECT().
				CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
					suite.sentCall = call
					return packedTextReturn(t, "https://node.example/"), nil
				}).
				Times(1)
		}

		t.Run("Should return it in a single call, without a trailing slash", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			initSuite(suite)

			url, err := suite.sut.ResolveUrl(resolvableEnsDomain)

			assert.NoError(t, err)
			assert.Equal(t, "https://node.example", url)
		})

		t.Run("Should call text(namehash, \"url\") on the configured registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveUrl(resolvableEnsDomain)
			assert.NoError(t, err)

			registryAbi, err := contract_abi.AgentEnsRegistryMetaData.GetAbi()
			assert.NoError(t, err)

			expectedQuery, err := registryAbi.Pack("text", helpers.Namehash(resolvableEnsDomain), "url")
			assert.NoError(t, err)

			if assert.NotNil(t, suite.sentCall.To) {
				assert.Equal(t, common.HexToAddress(agentEnsRegistryAddressForTest), *suite.sentCall.To)
			}
			assert.Equal(t, expectedQuery, suite.sentCall.Data)
		})
	})

	t.Run("Given the agent published no url", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingEnsUrlTestingSuite) {
			suite.mockEthClient.EXPECT().
				CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(packedTextReturn(t, ""), nil).
				Times(1)
		}

		t.Run("Should return an empty url rather than an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			initSuite(suite)

			url, err := suite.sut.ResolveUrl(resolvableEnsDomain)

			assert.NoError(t, err)
			assert.Empty(t, url)
		})
	})

	t.Run("Given the rpc call fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingEnsUrlTestingSuite) {
			suite.mockEthClient.EXPECT().
				CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveUrl(resolvableEnsDomain)

			assert.Error(t, err)
		})
	})

	// Reading another chain's name against this registry would answer empty, not wrong-but-obvious.
	t.Run("Given a name belonging to another chain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingEnsUrlTestingSuite) {
			suite.mockEthClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should refuse before calling the chain", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveUrl("agent.collection.1.6022.eth")

			assert.ErrorContains(t, err, "belongs to chain 1")
		})
	})

	t.Run("Given a name with no chain id label", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingEnsUrlTestingSuite) {
			suite.mockEthClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should refuse before calling the chain", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
		})
	})
}
