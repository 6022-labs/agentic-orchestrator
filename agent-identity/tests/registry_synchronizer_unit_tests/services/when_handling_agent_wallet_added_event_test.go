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

type WhenHandlingAgentWalletAddedEventTestingSuite struct {
	agentRepositoryMock                    *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionBlockchainRequesterMock *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.AgentWalletAddedEventHandler
}

func WhenHandlingAgentWalletAddedEventBeforeEach(t *testing.T) *WhenHandlingAgentWalletAddedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionBlockchainRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)

	sut := services.NewAgentWalletAddedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentCollectionBlockchainRequesterMock,
	)

	return &WhenHandlingAgentWalletAddedEventTestingSuite{
		agentRepositoryMock:                    agentRepositoryMock,
		agentCollectionBlockchainRequesterMock: agentCollectionBlockchainRequesterMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentWalletAddedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, wallet common.Address) *abi.AgentCollectionV1AgentWalletAdded {
		return &abi.AgentCollectionV1AgentWalletAdded{
			TokenId: new(big.Int).SetUint64(tokenID),
			Wallet:  wallet,
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given a wallet was added on-chain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentWalletAddedEventTestingSuite) (*abi.AgentCollectionV1AgentWalletAdded, common.Address) {
			collectionAddress := common_tests.RandomEthAddress()
			newWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, newWallet)

			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(agent, nil).
				Times(1)

			suite.agentCollectionBlockchainRequesterMock.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{Wallets: []common.Address{newWallet}}, nil).
				Times(1)

			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, newWallet.Hex(), saved.WalletAddress)
					return saved, nil
				}).
				Times(1)

			return event, newWallet
		}

		t.Run("Should re-sync the agent wallet from the chain", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentWalletAddedEventBeforeEach(t)
			event, _ := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given agent lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentWalletAddedEventTestingSuite) *abi.AgentCollectionV1AgentWalletAdded {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, common_tests.RandomEthAddress())

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress.Hex(), uint64(7)).
				Return(nil, assert.AnError).
				Times(1)

			return event
		}

		t.Run("Should return wrapped error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentWalletAddedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection")
		})
	})
}
