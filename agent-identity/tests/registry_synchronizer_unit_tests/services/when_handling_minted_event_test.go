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
	registry_common_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/samber/lo"
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
	decentralizedDomainNameBuilder := registry_common_interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(
		mockController,
	)
	agentCollectionBlockchainRequester := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(
		mockController,
	)

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

	// expectAgentInformation returns "Alice" with one evm address from the chain.
	expectAgentInformation := func(
		suite *WhenHandlingMintedEventTestingSuite,
		collectionAddress common.Address,
		wallet common.Address,
	) {
		suite.agentCollectionBlockchainRequester.
			EXPECT().
			GetAgentInformation(gomock.Any(), collectionAddress, gomock.Any()).
			Return(&abi.AgentInformation{
				Name:      "Alice",
				Creator:   common_tests.RandomEthAddress(),
				Addresses: []abi.AgentAddress{{AddressType: "evm", Value: wallet.Hex()}},
			}, nil)
	}

	// expectTokenRow returns the given row for the token's own lookup.
	expectTokenRow := func(
		suite *WhenHandlingMintedEventTestingSuite,
		collectionAddress common.Address,
		tokenID uint64,
		agent *models.Agent,
	) {
		suite.agentRepositoryMock.
			EXPECT().
			GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), tokenID).
			Return(agent, nil)
	}

	// expectProposal returns the given proposal for "Alice" holding the wallet.
	expectProposal := func(
		suite *WhenHandlingMintedEventTestingSuite,
		collectionAddress common.Address,
		wallet common.Address,
		proposal *models.Agent,
	) {
		suite.agentRepositoryMock.
			EXPECT().
			GetNotYetMintedByCollectionAddressNameAndAddresses(
				gomock.Any(),
				strings.ToLower(collectionAddress.Hex()),
				"Alice",
				map[enums.AddressType]string{enums.AddressTypeEvm: strings.ToLower(wallet.Hex())},
			).
			Return(proposal, nil)
	}

	// expectAgentSaved expects the rest of the sync for "Alice" with one evm address.
	expectAgentSaved := func(
		suite *WhenHandlingMintedEventTestingSuite,
		collectionAddress common.Address,
		wallet common.Address,
		assertSaved func(saved *models.Agent),
	) {
		owner := common_tests.RandomEthAddress()
		suite.agentCollectionBlockchainRequester.
			EXPECT().
			GetAgentOwner(gomock.Any(), collectionAddress, gomock.Any()).
			Return(&owner, nil)
		suite.agentCollectionRepositoryMock.
			EXPECT().
			GetByAddress(gomock.Any(), strings.ToLower(collectionAddress.Hex())).
			Return(model_builders.NewAgentCollectionBuilder().WithName("Genesis").Build(), nil)
		suite.decentralizedDomainNameBuilder.
			EXPECT().
			BuildDomainName("Alice", "Genesis").
			Return("alice.genesis.eth")
		suite.agentRepositoryMock.
			EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, saved *models.Agent) (*models.Agent, error) {
				assertSaved(saved)
				return saved, nil
			})
		suite.agentAddressRepositoryMock.
			EXPECT().
			ReplaceAllForAgent(gomock.Any(), gomock.Any(), map[enums.AddressType]string{
				enums.AddressTypeEvm: strings.ToLower(wallet.Hex()),
			}).
			Return(nil)
	}

	t.Run("Given fetching agent information fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(gomock.Any(), collectionAddress, gomock.Any()).
				Return(nil, assert.AnError)
			return newEvent(collectionAddress, 1)
		}

		t.Run("Should return wrapped info error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.ErrorIs(t, err, assert.AnError)
			assert.ErrorContains(t, err, "failed to get agent information")
		})
	})

	t.Run("Given the token lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			expectAgentInformation(suite, collectionAddress, common_tests.RandomEthAddress())
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(collectionAddress.Hex()), uint64(3)).
				Return(nil, assert.AnError)
			return newEvent(collectionAddress, 3)
		}

		t.Run("Should return wrapped token lookup error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.ErrorIs(t, err, assert.AnError)
			assert.ErrorContains(t, err, "failed to get agent by collection and token id")
		})
	})

	t.Run("Given the token is already synced", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) (*abi.AgentCollectionV1Minted, *models.Agent) {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			existing := model_builders.NewAgentBuilder().
				WithCollectionAddress(strings.ToLower(collectionAddress.Hex())).
				WithCollectionTokenId(lo.ToPtr(uint64(7))).
				Build()
			expectAgentInformation(suite, collectionAddress, wallet)
			expectTokenRow(suite, collectionAddress, 7, existing)
			expectAgentSaved(suite, collectionAddress, wallet, func(saved *models.Agent) {
				assert.Same(t, existing, saved)
			})
			return newEvent(collectionAddress, 7), existing
		}

		t.Run("Should update the token's own agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event, existing := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
			assert.Equal(t, "Alice", existing.Name)
		})
	})

	t.Run("Given the proposal lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			expectAgentInformation(suite, collectionAddress, common_tests.RandomEthAddress())
			expectTokenRow(suite, collectionAddress, 2, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressNameAndAddresses(
					gomock.Any(),
					strings.ToLower(collectionAddress.Hex()),
					"Alice",
					gomock.Any(),
				).
				Return(nil, assert.AnError)
			return newEvent(collectionAddress, 2)
		}

		t.Run("Should return wrapped proposal lookup error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.ErrorIs(t, err, assert.AnError)
			assert.ErrorContains(t, err, "failed to get proposal by name and addresses")
		})
	})

	t.Run("Given no proposal of the collection matches", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			expectAgentInformation(suite, collectionAddress, wallet)
			expectTokenRow(suite, collectionAddress, 1, nil)
			expectProposal(suite, collectionAddress, wallet, nil)
			expectAgentSaved(suite, collectionAddress, wallet, func(saved *models.Agent) {
				assert.Equal(t, strings.ToLower(collectionAddress.Hex()), saved.CollectionAddress)
				if assert.NotNil(t, saved.CollectionTokenId) {
					assert.Equal(t, uint64(1), *saved.CollectionTokenId)
				}
			})
			return newEvent(collectionAddress, 1)
		}

		t.Run("Should save a new agent for the token", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the mint carries only a btc address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:       "Agent",
					Attributes: []abi.KeyValue{{Key: "role", Value: "expert"}},
					Creator:    common_tests.RandomEthAddress(),
					Addresses: []abi.AgentAddress{
						{AddressType: "btc", Value: "bc1qexampleonly"},
					},
				}, nil)
			expectTokenRow(suite, collectionAddress, 6, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressNameAndAddresses(
					gomock.Any(),
					strings.ToLower(collectionAddress.Hex()),
					"Agent",
					map[enums.AddressType]string{enums.AddressTypeBtc: "bc1qexampleonly"},
				).
				Return(nil, nil)

			owner := common_tests.RandomEthAddress()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&owner, nil)

			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(gomock.Any(), strings.ToLower(collectionAddress.Hex())).
				Return(model_builders.NewAgentCollectionBuilder().WithName("Genesis").Build(), nil)

			suite.decentralizedDomainNameBuilder.
				EXPECT().
				BuildDomainName("Agent", "Genesis").
				Return("agent.genesis.eth")

			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, saved *models.Agent) (*models.Agent, error) {
					return saved, nil
				})

			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(gomock.Any(), gomock.Any(), map[enums.AddressType]string{
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil)

			return newEvent(collectionAddress, 6)
		}

		t.Run("Should sync the agent with the btc entry only", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given agent collection does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			expectAgentInformation(suite, collectionAddress, wallet)
			expectTokenRow(suite, collectionAddress, 4, nil)
			expectProposal(suite, collectionAddress, wallet, nil)

			owner := common_tests.RandomEthAddress()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&owner, nil)

			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(gomock.Any(), strings.ToLower(collectionAddress.Hex())).
				Return(nil, nil)

			return newEvent(collectionAddress, 4)
		}

		t.Run("Should return collection not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.ErrorContains(t, err, "agent collection with address")
		})
	})

	t.Run("Given the mint executes a proposal of the collection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			tokenID := uint64(5)
			wallet := common_tests.RandomEthAddress()
			creator := common_tests.RandomEthAddress()
			owner := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithName("Alice").
				WithCollectionAddress(strings.ToLower(collectionAddress.Hex())).
				WithCollectionTokenId(nil).
				WithImages(nil).
				WithAttributes(nil).
				Build()
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:    "Alice",
					Creator: creator,
					Addresses: []abi.AgentAddress{
						{AddressType: "btc", Value: "bc1qexampleonly"},
						{AddressType: "evm", Value: wallet.Hex()},
					},
					Images: []abi.KeyValue{{Key: "avatar", Value: "ipfs://avatar"}},
					Attributes: []abi.KeyValue{
						{Key: "role", Value: "human"},
						{Key: "skill", Value: "analysis"},
					},
				}, nil)
			expectTokenRow(suite, collectionAddress, tokenID, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressNameAndAddresses(
					gomock.Any(),
					strings.ToLower(collectionAddress.Hex()),
					"Alice",
					map[enums.AddressType]string{
						enums.AddressTypeEvm: strings.ToLower(wallet.Hex()),
						enums.AddressTypeBtc: "bc1qexampleonly",
					},
				).
				Return(agent, nil)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(gomock.Any(), collectionAddress, gomock.Any()).
				Return(&owner, nil)
			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(gomock.Any(), strings.ToLower(collectionAddress.Hex())).
				Return(model_builders.NewAgentCollectionBuilder().WithName("Genesis").Build(), nil)
			suite.decentralizedDomainNameBuilder.
				EXPECT().
				BuildDomainName("Alice", "Genesis").
				Return("alice.genesis.eth")
			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any(), agent).
				DoAndReturn(func(_ context.Context, saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, "Alice", saved.Name)
					assert.Equal(t, "alice.genesis.eth", saved.EnsDomain)
					assert.Equal(t, strings.ToLower(collectionAddress.Hex()), saved.CollectionAddress)
					if assert.NotNil(t, saved.CollectionTokenId) {
						assert.Equal(t, tokenID, *saved.CollectionTokenId)
					}
					assert.Nil(t, saved.ProposalId)
					assert.Equal(t, "ipfs://avatar", saved.Images["avatar"])
					assert.Equal(t, "analysis", saved.Attributes["skill"])
					assert.Equal(t, "human", saved.Attributes["role"])
					if assert.NotNil(t, saved.Owner) {
						assert.Equal(t, strings.ToLower(owner.Hex()), *saved.Owner)
					}
					assert.Equal(t, strings.ToLower(creator.Hex()), saved.Creator)
					return saved, nil
				})
			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(gomock.Any(), agent.ID, map[enums.AddressType]string{
					enums.AddressTypeEvm: strings.ToLower(wallet.Hex()),
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil)
			return newEvent(collectionAddress, tokenID)
		}

		t.Run("Should update and save the proposal agent from blockchain information", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})
}
