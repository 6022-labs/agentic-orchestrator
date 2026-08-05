package services_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingAgentAttributeUpdatedEventTestingSuite struct {
	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface

	sut *services.AgentAttributeUpdatedEventHandler
}

func WhenHandlingAgentAttributeUpdatedEventBeforeEach(t *testing.T) *WhenHandlingAgentAttributeUpdatedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := services.NewAgentAttributeUpdatedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
	)

	return &WhenHandlingAgentAttributeUpdatedEventTestingSuite{
		agentRepositoryMock: agentRepositoryMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentAttributeUpdatedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, key, oldValue, newValue string) *abi.AgentCollectionV1AgentAttributeUpdated {
		return &abi.AgentCollectionV1AgentAttributeUpdated{
			TokenId:  new(big.Int).SetUint64(tokenID),
			Key:      common.HexToHash(key),
			OldValue: oldValue,
			NewValue: newValue,
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given agent lookup fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAttributeUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentAttributeUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, "0x01", "old", "new")

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAttributeUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection")
		})
	})

	t.Run("Given agent does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAttributeUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentAttributeUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 8, "0x02", "old", "new")

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(8)).
				Return(nil, nil)

			return event
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAttributeUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "agent not found")
		})
	})

	t.Run("Given saving agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAttributeUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentAttributeUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 9, "0x03", "old", "new")
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(collectionAddress.Hex()).
				WithCollectionTokenId(lo.ToPtr(uint64(9))).
				WithAttributes(map[string]string{}).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(9)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAttributeUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to update agent")
		})
	})

	t.Run("Given agent is found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAttributeUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentAttributeUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 10, "0x04", "old", "new-value")
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(collectionAddress.Hex()).
				WithCollectionTokenId(lo.ToPtr(uint64(10))).
				WithAttributes(nil).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(10)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, "new-value", saved.Attributes[event.Key.String()])
					return saved, nil
				})

			return event
		}

		t.Run("Should upsert attribute and save agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAttributeUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
