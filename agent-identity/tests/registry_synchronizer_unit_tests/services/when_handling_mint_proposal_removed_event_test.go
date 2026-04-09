package services_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
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

type WhenHandlingMintProposalRemovedEventTestingSuite struct {
	agentRepositoryMock           *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRepositoryMock *repositories_mocks.MockAgentCollectionRepositoryInterface

	sut *services.MintProposalRemovedEventHandler
}

func WhenHandlingMintProposalRemovedEventBeforeEach(t *testing.T) *WhenHandlingMintProposalRemovedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	sut := services.NewMintProposalRemovedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentCollectionRepositoryMock,
	)

	return &WhenHandlingMintProposalRemovedEventTestingSuite{
		agentRepositoryMock:           agentRepositoryMock,
		agentCollectionRepositoryMock: agentCollectionRepositoryMock,

		sut: sut,
	}
}

func TestWhenHandlingMintProposalRemovedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, proposalID uint64) *abi.AgentCollectionV1MintProposalRemoved {
		return &abi.AgentCollectionV1MintProposalRemoved{
			Owner:      common_tests.RandomEthAddress(),
			Signer:     common_tests.RandomEthAddress(),
			ProposalId: new(big.Int).SetUint64(proposalID),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given agent lookup fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalRemovedEventTestingSuite) *abi.AgentCollectionV1MintProposalRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 1)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(collectionAddress.Hex(), uint64(1)).
				Return(nil, assert.AnError)
			return event
		}

		t.Run("Should return repository error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalRemovedEventTestingSuite) *abi.AgentCollectionV1MintProposalRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 2)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(collectionAddress.Hex(), uint64(2)).
				Return(nil, gorm.ErrRecordNotFound)
			return event
		}

		t.Run("Should return not found message", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "no agent found")
		})
	})

	t.Run("Given removing agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalRemovedEventTestingSuite) *abi.AgentCollectionV1MintProposalRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 3)
			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(collectionAddress.Hex(), uint64(3)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Remove(agent).
				Return(assert.AnError)

			return event
		}

		t.Run("Should return wrapped remove error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to remove agent from proposal")
		})
	})

	t.Run("Given event is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalRemovedEventTestingSuite) *abi.AgentCollectionV1MintProposalRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 4)
			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(collectionAddress.Hex(), uint64(4)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Remove(agent).
				Return(nil)

			return event
		}

		t.Run("Should remove matched agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
