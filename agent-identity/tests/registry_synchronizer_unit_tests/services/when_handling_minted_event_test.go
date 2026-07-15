package services_test

import (
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
	registry_common_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingMintedEventTestingSuite struct {
	agentRepositoryMock                *repositories_mocks.MockAgentRepositoryInterface
	agentAddressRepositoryMock         *repositories_mocks.MockAgentAddressRepositoryInterface
	agentCollectionRepositoryMock      *repositories_mocks.MockAgentCollectionRepositoryInterface
	decentralizedDomainNameBuilder     *registry_common_interfaces_mocks.MockDecentralizedDomainNameBuilderInterface
	agentCollectionBlockchainRequester *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.MintedEventHandler
}

func WhenHandlingMintedEventBeforeEach(t *testing.T) *WhenHandlingMintedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentAddressRepositoryMock := repositories_mocks.NewMockAgentAddressRepositoryInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	decentralizedDomainNameBuilder := registry_common_interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(mockController)
	agentCollectionBlockchainRequester := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)

	sut := services.NewMintedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentAddressRepositoryMock,
		agentCollectionRepositoryMock,
		agentCollectionBlockchainRequester,
		decentralizedDomainNameBuilder,
	)

	return &WhenHandlingMintedEventTestingSuite{
		agentRepositoryMock:                agentRepositoryMock,
		agentAddressRepositoryMock:         agentAddressRepositoryMock,
		agentCollectionRepositoryMock:      agentCollectionRepositoryMock,
		decentralizedDomainNameBuilder:     decentralizedDomainNameBuilder,
		agentCollectionBlockchainRequester: agentCollectionBlockchainRequester,

		sut: sut,
	}
}

func TestWhenHandlingMintedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, tokenID uint64) *abi.AgentCollectionV1Minted {
		return &abi.AgentCollectionV1Minted{
			To:      common_tests.RandomEthAddress(),
			TokenId: new(big.Int).SetUint64(tokenID),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given fetching agent information fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 1)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(nil, assert.AnError)
			return event
		}

		t.Run("Should return wrapped info error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent information")
		})
	})

	t.Run("Given agent address lookup fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 2)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:    "Agent",
					Role:    "expert",
					Creator: common_tests.RandomEthAddress(),
					Addresses: []abi.AgentAddress{
						{AddressType: "btc", Value: "bc1qexampleonly"},
						{AddressType: "evm", Value: strings.ToLower(wallet.Hex())},
					},
				}, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(wallet.Hex())).
				Return(nil, assert.AnError)
			return event
		}

		t.Run("Should return wrapped address lookup error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to find agent by addresses")
		})
	})

	t.Run("Given agent collection does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 4)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:    "Agent",
					Role:    "expert",
					Creator: common_tests.RandomEthAddress(),
					Addresses: []abi.AgentAddress{
						{AddressType: "btc", Value: "bc1qexampleonly"},
						{AddressType: "evm", Value: strings.ToLower(wallet.Hex())},
					},
				}, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(wallet.Hex())).
				Return(nil, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeBtc, "bc1qexampleonly").
				Return(nil, nil)

			owner := common_tests.RandomEthAddress()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(collectionAddress, gomock.Any()).
				Return(&owner, nil)

			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(collectionAddress.Hex()).
				Return(nil, nil)

			return event
		}

		t.Run("Should return collection not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "agent collection with address")
		})
	})

	t.Run("Given the mint carries only a btc address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 6)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:    "Agent",
					Role:    "expert",
					Creator: common_tests.RandomEthAddress(),
					Addresses: []abi.AgentAddress{
						{AddressType: "btc", Value: "bc1qexampleonly"},
					},
				}, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeBtc, "bc1qexampleonly").
				Return(nil, nil)

			owner := common_tests.RandomEthAddress()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(collectionAddress, gomock.Any()).
				Return(&owner, nil)

			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(collectionAddress.Hex()).
				Return(model_builders.NewAgentCollectionBuilder().WithName("Genesis").Build(), nil)

			suite.decentralizedDomainNameBuilder.
				EXPECT().
				BuildDomainName("Agent", "Genesis").
				Return("agent.genesis.eth")

			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					return saved, nil
				})

			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(gomock.Any(), map[enums.AddressType]string{
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil)

			return event
		}

		t.Run("Should sync the agent with the btc entry only", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given event is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			tokenID := uint64(5)
			wallet := common_tests.RandomEthAddress()
			owner := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithImages(nil).
				WithAttributes(nil).
				Build()
			event := newEvent(collectionAddress, tokenID)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:    "Alice",
					Role:    "human",
					Creator: common_tests.RandomEthAddress(),
					Addresses: []abi.AgentAddress{
						{AddressType: "btc", Value: "bc1qexampleonly"},
						{AddressType: "evm", Value: strings.ToLower(wallet.Hex())},
					},
					Images:     []abi.KeyValue{{Key: "avatar", Value: "ipfs://avatar"}},
					Attributes: []abi.KeyValue{{Key: "skill", Value: "analysis"}},
				}, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(wallet.Hex())).
				Return(model_builders.NewAgentAddressBuilder().WithAgentId(agent.ID).WithAddress(strings.ToLower(wallet.Hex())).Build(), nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetById(agent.ID).
				Return(agent, nil)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(collectionAddress, gomock.Any()).
				Return(&owner, nil)
			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(collectionAddress.Hex()).
				Return(model_builders.NewAgentCollectionBuilder().WithName("Genesis").Build(), nil)
			suite.decentralizedDomainNameBuilder.
				EXPECT().
				BuildDomainName("Alice", "Genesis").
				Return("alice.genesis.eth")
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, "Alice", saved.Name)
					assert.Equal(t, "alice.genesis.eth", saved.EnsDomain)
					assert.Equal(t, collectionAddress.Hex(), saved.CollectionAddress)
					if assert.NotNil(t, saved.CollectionTokenId) {
						assert.Equal(t, tokenID, *saved.CollectionTokenId)
					}
					assert.Nil(t, saved.ProposalId)
					assert.Equal(t, "ipfs://avatar", saved.Images["avatar"])
					assert.Equal(t, "analysis", saved.Attributes["skill"])
					if assert.NotNil(t, saved.Owner) {
						assert.Equal(t, owner.Hex(), *saved.Owner)
					}
					return saved, nil
				})
			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(agent.ID, map[enums.AddressType]string{
					enums.AddressTypeEvm: strings.ToLower(wallet.Hex()),
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil)
			return event
		}

		t.Run("Should update and save agent from blockchain information", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
