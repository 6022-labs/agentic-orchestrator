package services_test

import (
	"io"
	"log"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenHandlingAgentSignerUpdatedEventTestingSuite struct {
	agentRepositoryMock          *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRequesterMock *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface
	orchestratorPublicKey        common.Address

	sut *services.AgentSignerUpdatedEventHandler
}

func WhenHandlingAgentSignerUpdatedEventBeforeEach(t *testing.T) *WhenHandlingAgentSignerUpdatedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)

	orchestratorPublicKey := common_tests.RandomEthAddress()

	sut := services.NewAgentSignerUpdatedEventHandler(
		log.New(io.Discard, "", 0),
		settings.OrchestratorSettings{PublicKey: orchestratorPublicKey},
		agentRepositoryMock,
		agentCollectionRequesterMock,
	)

	return &WhenHandlingAgentSignerUpdatedEventTestingSuite{
		agentRepositoryMock:          agentRepositoryMock,
		agentCollectionRequesterMock: agentCollectionRequesterMock,
		orchestratorPublicKey:        orchestratorPublicKey,

		sut: sut,
	}
}

func TestWhenHandlingAgentSignerUpdatedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, newSigner common.Address) *abi.AgentCollectionV1AgentSignerUpdated {
		return &abi.AgentCollectionV1AgentSignerUpdated{
			TokenId:   new(big.Int).SetUint64(tokenID),
			OldSigner: common_tests.RandomEthAddress(),
			NewSigner: newSigner,
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given agent lookup returns unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentSignerUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentSignerUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 1, common_tests.RandomEthAddress())
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(1)).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped get error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentSignerUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection and id")
		})
	})

	t.Run("Given signer is external", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentSignerUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentSignerUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 2, common_tests.RandomEthAddress())
			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(2)).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.False(t, saved.IsInternal)
					return saved, nil
				})

			return event
		}

		t.Run("Should mark agent as external and save", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentSignerUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given signer is internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentSignerUpdatedEventTestingSuite) *abi.AgentCollectionV1AgentSignerUpdated {
			collectionAddress := common_tests.RandomEthAddress()
			tokenID := uint64(3)
			event := newEvent(collectionAddress, tokenID, suite.orchestratorPublicKey)
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&tokenID).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), tokenID).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.True(t, saved.IsInternal)
					return saved, nil
				})

			return event
		}

		t.Run("Should mark agent as internal and save", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentSignerUpdatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
