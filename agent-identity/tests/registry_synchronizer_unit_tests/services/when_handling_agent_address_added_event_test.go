package services_test

import (
	"context"

	"math/big"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
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

type WhenHandlingAgentAddressAddedEventTestingSuite struct {
	agentRepositoryMock                    *repositories_mocks.MockAgentRepositoryInterface
	agentAddressRepositoryMock             *repositories_mocks.MockAgentAddressRepositoryInterface
	agentCollectionBlockchainRequesterMock *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.AgentAddressAddedEventHandler
}

func WhenHandlingAgentAddressAddedEventBeforeEach(t *testing.T) *WhenHandlingAgentAddressAddedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentAddressRepositoryMock := repositories_mocks.NewMockAgentAddressRepositoryInterface(mockController)
	agentCollectionBlockchainRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(
		mockController,
	)

	sut := services.NewAgentAddressAddedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentAddressRepositoryMock,
		agentCollectionBlockchainRequesterMock,
	)

	return &WhenHandlingAgentAddressAddedEventTestingSuite{
		agentRepositoryMock:                    agentRepositoryMock,
		agentAddressRepositoryMock:             agentAddressRepositoryMock,
		agentCollectionBlockchainRequesterMock: agentCollectionBlockchainRequesterMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentAddressAddedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, wallet common.Address) *abi.AgentCollectionV1AgentAddressAdded {
		return &abi.AgentCollectionV1AgentAddressAdded{
			TokenId:     new(big.Int).SetUint64(tokenID),
			AddressType: "evm",
			Value:       strings.ToLower(wallet.Hex()),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given the event carries an evm value that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should skip it without error and without touching the repositories", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAddressAddedEventBeforeEach(t)
			event := newEvent(common_tests.RandomEthAddress(), 7, common_tests.RandomEthAddress())
			event.Value = "0xnothex"

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given agent lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAddressAddedEventTestingSuite) *abi.AgentCollectionV1AgentAddressAdded {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, common_tests.RandomEthAddress())

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(7)).
				Return(nil, assert.AnError).
				Times(1)

			return event
		}

		t.Run("Should return wrapped error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAddressAddedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection")
		})
	})

	t.Run("Given the chain addresses already match the mirror", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAddressAddedEventTestingSuite) *abi.AgentCollectionV1AgentAddressAdded {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, wallet)

			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(7)).
				Return(agent, nil).
				Times(1)

			suite.agentCollectionBlockchainRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{Addresses: []abi.AgentAddress{
					{AddressType: "btc", Value: "bc1qexampleonly"},
					{AddressType: "evm", Value: strings.ToLower(wallet.Hex())},
				}}, nil).
				Times(1)

			suite.agentAddressRepositoryMock.
				EXPECT().
				GetAllByAgentId(gomock.Any(), agent.ID).
				Return([]*models.AgentAddress{
					model_builders.NewAgentAddressBuilder().
						WithAgentId(agent.ID).
						WithAddress(strings.ToLower(wallet.Hex())).
						Build(),
					model_builders.NewAgentAddressBuilder().
						WithAgentId(agent.ID).
						WithType(enums.AddressTypeBtc).
						WithAddress("bc1qexampleonly").
						Build(),
				}, nil).
				Times(1)

			return event
		}

		t.Run("Should not replace the mirrored addresses", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAddressAddedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an address was added on-chain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAddressAddedEventTestingSuite) (*abi.AgentCollectionV1AgentAddressAdded, common.Address) {
			collectionAddress := common_tests.RandomEthAddress()
			newWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 7, newWallet)

			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(7)).
				Return(agent, nil).
				Times(1)

			suite.agentCollectionBlockchainRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{Addresses: []abi.AgentAddress{
					{AddressType: "evm", Value: "not-an-address"},
					{AddressType: "btc", Value: "BC1QW508D6QEJXTDG4Y5R3ZARVARY0C5XW7KV8F3T4"},
					{AddressType: "evm", Value: newWallet.Hex()},
				}}, nil).
				Times(1)

			suite.agentAddressRepositoryMock.
				EXPECT().
				GetAllByAgentId(gomock.Any(), agent.ID).
				Return([]*models.AgentAddress{
					model_builders.NewAgentAddressBuilder().WithAgentId(agent.ID).WithAddress("0xstalevalue").Build(),
				}, nil).
				Times(1)

			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(gomock.Any(), agent.ID, map[enums.AddressType]string{
					enums.AddressTypeEvm: strings.ToLower(newWallet.Hex()),
					enums.AddressTypeBtc: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
				}).
				Return(nil).
				Times(1)

			return event, newWallet
		}

		t.Run("Should re-sync the valid addresses, lowercasing evm and btc", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAddressAddedEventBeforeEach(t)
			event, _ := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})
}
