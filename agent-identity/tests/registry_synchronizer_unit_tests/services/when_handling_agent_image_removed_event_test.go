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
)

type WhenHandlingAgentImageRemovedEventTestingSuite struct {
	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface

	sut *services.AgentImageRemovedEventHandler
}

func WhenHandlingAgentImageRemovedEventBeforeEach(t *testing.T) *WhenHandlingAgentImageRemovedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := services.NewAgentImageRemovedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
	)

	return &WhenHandlingAgentImageRemovedEventTestingSuite{
		agentRepositoryMock: agentRepositoryMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentImageRemovedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, key string) *abi.AgentCollectionV1AgentImageRemoved {
		return &abi.AgentCollectionV1AgentImageRemoved{
			TokenId: new(big.Int).SetUint64(tokenID),
			Key:     common.HexToHash(key),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given agent lookup fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageRemovedEventTestingSuite) *abi.AgentCollectionV1AgentImageRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 4, "0x01")
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(4)).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped get error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection")
		})
	})

	t.Run("Given agent does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageRemovedEventTestingSuite) *abi.AgentCollectionV1AgentImageRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 5, "0x02")
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(5)).
				Return(nil, nil)

			return event
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "agent not found")
		})
	})

	t.Run("Given agent has no images", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageRemovedEventTestingSuite) *abi.AgentCollectionV1AgentImageRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 6, "0x03")
			agent := model_builders.NewAgentBuilder().
				WithImages(nil).
				Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(6)).
				Return(agent, nil)

			return event
		}

		t.Run("Should do nothing and not persist", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given saving agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageRemovedEventTestingSuite) *abi.AgentCollectionV1AgentImageRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, "0x04")
			agent := model_builders.NewAgentBuilder().
				WithImages(map[string]string{event.Key.String(): "value"}).
				Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to update agent")
		})
	})

	t.Run("Given agent is found with the image", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentImageRemovedEventTestingSuite) *abi.AgentCollectionV1AgentImageRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 8, "0x05")
			agent := model_builders.NewAgentBuilder().
				WithImages(map[string]string{
					event.Key.String(): "removed-image",
					"other":            "kept-image",
				}).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(8)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.NotContains(t, saved.Images, event.Key.String())
					assert.Equal(t, "kept-image", saved.Images["other"])
					return saved, nil
				})

			return event
		}

		t.Run("Should remove the image and persist", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentImageRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
