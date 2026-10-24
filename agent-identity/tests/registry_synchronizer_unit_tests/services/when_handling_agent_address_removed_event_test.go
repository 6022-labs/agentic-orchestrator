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

type WhenHandlingAgentAddressRemovedEventTestingSuite struct {
	agentRepositoryMock                    *repositories_mocks.MockAgentRepositoryInterface
	agentAddressRepositoryMock             *repositories_mocks.MockAgentAddressRepositoryInterface
	agentCollectionBlockchainRequesterMock *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.AgentAddressRemovedEventHandler
}

func WhenHandlingAgentAddressRemovedEventBeforeEach(t *testing.T) *WhenHandlingAgentAddressRemovedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentAddressRepositoryMock := repositories_mocks.NewMockAgentAddressRepositoryInterface(mockController)
	agentCollectionBlockchainRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(
		mockController,
	)

	sut := services.NewAgentAddressRemovedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentAddressRepositoryMock,
		agentCollectionBlockchainRequesterMock,
	)

	return &WhenHandlingAgentAddressRemovedEventTestingSuite{
		agentRepositoryMock:                    agentRepositoryMock,
		agentAddressRepositoryMock:             agentAddressRepositoryMock,
		agentCollectionBlockchainRequesterMock: agentCollectionBlockchainRequesterMock,

		sut: sut,
	}
}

func TestWhenHandlingAgentAddressRemovedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64, wallet common.Address) *abi.AgentCollectionV1AgentAddressRemoved {
		return &abi.AgentCollectionV1AgentAddressRemoved{
			TokenId:     new(big.Int).SetUint64(tokenID),
			AddressType: "evm",
			Value:       strings.ToLower(wallet.Hex()),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given the event carries an evm value without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should skip it without error and without touching the repositories", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAddressRemovedEventBeforeEach(t)
			event := newEvent(common_tests.RandomEthAddress(), 7, common_tests.RandomEthAddress())
			event.Value = strings.TrimPrefix(event.Value, "0x")

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given agent lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAddressRemovedEventTestingSuite) *abi.AgentCollectionV1AgentAddressRemoved {
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

			suite := WhenHandlingAgentAddressRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by collection")
		})
	})

	t.Run("Given the chain addresses already match the mirror", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAddressRemovedEventTestingSuite) *abi.AgentCollectionV1AgentAddressRemoved {
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

			suite := WhenHandlingAgentAddressRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an address was removed on-chain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingAgentAddressRemovedEventTestingSuite) *abi.AgentCollectionV1AgentAddressRemoved {
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

			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(gomock.Any(), agent.ID, map[enums.AddressType]string{
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil).
				Times(1)

			return event
		}

		t.Run("Should re-sync the agent addresses from the chain", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingAgentAddressRemovedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})
}
