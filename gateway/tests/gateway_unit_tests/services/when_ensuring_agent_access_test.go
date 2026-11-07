package services_test

import (
	"context"
	"strings"
	"testing"

	gateway_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenEnsuringAgentAccessTestingSuite struct {
	sut *services.AgentAccessVerifier

	mockAgentRegistryRequester *interfaces_mocks.MockAgentRegistryRequesterInterface
}

func WhenEnsuringAgentAccessBeforeEach(t *testing.T) *WhenEnsuringAgentAccessTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)

	return &WhenEnsuringAgentAccessTestingSuite{
		sut:                        services.NewAgentAccessVerifier(mockAgentRegistryRequester),
		mockAgentRegistryRequester: mockAgentRegistryRequester,
	}
}

func TestWhenEnsuringAgentAccess(t *testing.T) {
	t.Parallel()

	t.Run("Given the registry fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an unavailable error carrying the cause", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				collectionAddress,
				7,
			)

			var unavailable *gateway_errors.UnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "agent_registry_unavailable", unavailable.Code)
			assert.Equal(t, "agent registry is unavailable", unavailable.Message)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given the agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				collectionAddress,
				7,
			)

			var notFoundErr *gateway_errors.NotFoundError
			assert.ErrorAs(t, err, &notFoundErr)
			assert.Equal(t, "agent_not_found", notFoundErr.Code)
		})
	})

	t.Run("Given an agent without owner nor evm address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, collectionAddress string) {
			agent := model_builders.NewAgentBuilder().WithAddresses(map[string]string{}).Build()
			agent.Owner = nil

			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(agent, nil).
				Times(1)
		}

		t.Run("Should return a forbidden error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				collectionAddress,
				7,
			)

			var forbiddenErr *gateway_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.AgentAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given the caller has no access to the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(model_builders.NewAgentBuilder().WithOwner(common_tests.RandomEthAddress().Hex()).Build(), nil).
				Times(1)
		}

		t.Run("Should return a forbidden error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				collectionAddress,
				7,
			)

			var forbiddenErr *gateway_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.AgentAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, callerId string, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), strings.ToLower(collectionAddress), uint64(7)).
				Return(model_builders.NewAgentBuilder().WithOwner(callerId).Build(), nil).
				Times(1)
		}

		t.Run("Should look the agent up by its lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			collectionAddress := common_tests.RandomEthAddress().Hex()
			initSuite(suite, callerId, collectionAddress)

			err := suite.sut.EnsureAccess(context.Background(), callerId, collectionAddress, 7)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the caller has access to the agent as its owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, callerId string, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(model_builders.NewAgentBuilder().WithOwner(callerId).Build(), nil).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, callerId, collectionAddress)

			err := suite.sut.EnsureAccess(context.Background(), callerId, collectionAddress, 7)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an owner stored in checksum form", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, owner string, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(model_builders.NewAgentBuilder().WithOwner(owner).Build(), nil).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			owner := common_tests.RandomEthAddress().Hex()
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, owner, collectionAddress)

			err := suite.sut.EnsureAccess(context.Background(), strings.ToLower(owner), collectionAddress, 7)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the caller is one of the agent's evm addresses", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringAgentAccessTestingSuite, agentEvmAddress string, collectionAddress string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), collectionAddress, uint64(7)).
				Return(
					model_builders.NewAgentBuilder().
						WithOwner(common_tests.RandomEthAddress().Hex()).
						WithAddresses(map[string]string{"evm": agentEvmAddress}).
						Build(),
					nil,
				).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringAgentAccessBeforeEach(t)
			agentEvmAddress := common_tests.RandomEthAddress().Hex()
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, agentEvmAddress, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(agentEvmAddress),
				collectionAddress,
				7,
			)

			assert.NoError(t, err)
		})
	})
}
