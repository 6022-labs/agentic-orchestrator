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
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingAgentWalletRemovedEventTestingSuite struct {
	agentRepositoryMock                    *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionBlockchainRequesterMock *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.AgentWalletRemovedEventHandler
}

func WhenHandlingAgentWalletRemovedEventBeforeEach(t *testing.T) *WhenHandlingAgentWalletRemovedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionBlockchainRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)

	sut := services.NewAgentWalletRemovedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentCollectionBlockchainRequesterMock,
	)

	return &WhenHandlingAgentWalletRemovedEventTestingSuite{
		agentRepositoryMock:                    agentRepositoryMock,
		agentCollectionBlockchainRequesterMock: agentCollectionBlockchainRequesterMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentWalletRemovedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, wallet common.Address) *abi.AgentCollectionV1AgentWalletRemoved {
		return &abi.AgentCollectionV1AgentWalletRemoved{
			TokenId: new(big.Int).SetUint64(tokenID),
			Wallet:  wallet,
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given the chain wallet already matches the mirror", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentWalletRemovedEventTestingSuite) *abi.AgentCollectionV1AgentWalletRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, wallet)

			agent := model_builders.NewAgentBuilder().Build()
			agent.WalletAddress = wallet.Hex()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(agent, nil).
				Times(1)

			suite.agentCollectionBlockchainRequesterMock.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{Wallets: []common.Address{wallet}}, nil).
				Times(1)

			return event
		}

		t.Run("Should not save the agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentWalletRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the last wallet was removed on-chain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentWalletRemovedEventTestingSuite) *abi.AgentCollectionV1AgentWalletRemoved {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, wallet)

			agent := model_builders.NewAgentBuilder().Build()
			agent.WalletAddress = wallet.Hex()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(agent, nil).
				Times(1)

			suite.agentCollectionBlockchainRequesterMock.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{Wallets: []common.Address{}}, nil).
				Times(1)

			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, "", saved.WalletAddress)
					return saved, nil
				}).
				Times(1)

			return event
		}

		t.Run("Should clear the mirrored wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentWalletRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
