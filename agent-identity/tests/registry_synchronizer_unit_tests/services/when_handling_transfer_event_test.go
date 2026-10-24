package services_test

import (
	"context"
	"strings"

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

type WhenHandlingTransferEventTestingSuite struct {
	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface

	sut *services.TransferEventHandler
}

func WhenHandlingTransferEventBeforeEach(t *testing.T) *WhenHandlingTransferEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := services.NewTransferEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
	)

	return &WhenHandlingTransferEventTestingSuite{
		agentRepositoryMock: agentRepositoryMock,

		sut: sut,
	}
}

func TestWhenHandlingTransferEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, to common.Address) *abi.AgentCollectionV1Transfer {
		return &abi.AgentCollectionV1Transfer{
			From:    common_tests.RandomEthAddress(),
			To:      to,
			TokenId: new(big.Int).SetUint64(tokenID),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given agent does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingTransferEventTestingSuite) *abi.AgentCollectionV1Transfer {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 1, common_tests.RandomEthAddress())
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(1)).
				Return(nil, nil)

			return event
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingTransferEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "agent not found for transfer event")
		})
	})

	t.Run("Given agent lookup returns no agent and no error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingTransferEventTestingSuite) *abi.AgentCollectionV1Transfer {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 2, common_tests.RandomEthAddress())
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(2)).
				Return(nil, nil)

			return event
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingTransferEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "agent not found for transfer event")
		})
	})

	t.Run("Given saving agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingTransferEventTestingSuite) *abi.AgentCollectionV1Transfer {
			collectionAddress := common_tests.RandomEthAddress()
			newOwner := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 3, newOwner)
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(3)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any(), agent).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingTransferEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given event is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingTransferEventTestingSuite) *abi.AgentCollectionV1Transfer {
			collectionAddress := common_tests.RandomEthAddress()
			newOwner := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 4, newOwner)
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(4)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any(), agent).
				DoAndReturn(func(_ context.Context, saved *models.Agent) (*models.Agent, error) {
					if assert.NotNil(t, saved.Owner) {
						assert.Equal(t, strings.ToLower(newOwner.Hex()), *saved.Owner)
					}
					return saved, nil
				})

			return event
		}

		t.Run("Should update owner and save", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingTransferEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})
}
