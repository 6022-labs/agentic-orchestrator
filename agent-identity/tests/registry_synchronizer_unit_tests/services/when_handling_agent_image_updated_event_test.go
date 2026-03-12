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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenHandlingAgentImageUpdatedEventTestingSuite struct {
	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface

	sut *services.AgentImageUpdatedEventHandler
}

func WhenHandlingAgentImageUpdatedEventBeforeEach(t *testing.T) *WhenHandlingAgentImageUpdatedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := services.NewAgentImageUpdatedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
	)

	return &WhenHandlingAgentImageUpdatedEventTestingSuite{
		agentRepositoryMock: agentRepositoryMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentImageUpdatedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, key, oldValue, newValue string) *abi.AgentCollectionV1AgentImageUpdated {
		return &abi.AgentCollectionV1AgentImageUpdated{
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

		initSuite := func(suite *WhenHandlingAgentImageUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentImageUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 4, "0x01", "old", "new")
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(4)).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped get error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection and id")
		})
	})

	t.Run("Given agent does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentImageUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 5, "0x02", "old", "new")
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(5)).
				Return(nil, gorm.ErrRecordNotFound)

			return event
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "agent not found")
		})
	})

	t.Run("Given saving agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentImageUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 6, "0x03", "old", "new")
			agent := model_builders.NewAgentBuilder().
				WithImages(nil).
				Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(6)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to update agent")
		})
	})

	t.Run("Given agent is found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentImageUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, "0x04", "old", "new-image")
			agent := model_builders.NewAgentBuilder().
				WithImages(nil).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, "new-image", saved.Images[event.Key.String()])
					return saved, nil
				})

			return event
		}

		t.Run("Should set image and persist", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
