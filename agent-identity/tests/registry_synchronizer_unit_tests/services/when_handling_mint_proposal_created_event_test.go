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

type WhenHandlingMintProposalCreatedEventTestingSuite struct {
	agentRepositoryMock                *repositories_mocks.MockAgentRepositoryInterface
	agentAddressRepositoryMock         *repositories_mocks.MockAgentAddressRepositoryInterface
	agentCollectionRepositoryMock      *repositories_mocks.MockAgentCollectionRepositoryInterface
	agentCollectionRequesterMock       *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface
	decentralizedDomainNameBuilderMock *registry_common_interfaces_mocks.MockDecentralizedDomainNameBuilderInterface

	sut *services.MintProposalCreatedEventHandler
}

func WhenHandlingMintProposalCreatedEventBeforeEach(t *testing.T) *WhenHandlingMintProposalCreatedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentAddressRepositoryMock := repositories_mocks.NewMockAgentAddressRepositoryInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)
	decentralizedDomainNameBuilderMock := registry_common_interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(mockController)

	sut := services.NewMintProposalCreatedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentAddressRepositoryMock,
		agentCollectionRepositoryMock,
		agentCollectionRequesterMock,
		decentralizedDomainNameBuilderMock,
	)

	return &WhenHandlingMintProposalCreatedEventTestingSuite{
		agentRepositoryMock:                agentRepositoryMock,
		agentAddressRepositoryMock:         agentAddressRepositoryMock,
		agentCollectionRepositoryMock:      agentCollectionRepositoryMock,
		agentCollectionRequesterMock:       agentCollectionRequesterMock,
		decentralizedDomainNameBuilderMock: decentralizedDomainNameBuilderMock,

		sut: sut,
	}
}

func TestWhenHandlingMintProposalCreatedEvent(t *testing.T) {
	t.Parallel()

	// Proposals carry a btc entry first and lowercase evm values: handlers must skip
	// non-evm entries and normalize to the checksummed form.
	newMintProposal := func(proposalID uint64, wallets []common.Address) *abi.MintProposal {
		addresses := []abi.AgentAddress{{AddressType: "btc", Value: "bc1qexampleonly"}}
		for _, wallet := range wallets {
			addresses = append(addresses, abi.AgentAddress{AddressType: "evm", Value: strings.ToLower(wallet.Hex())})
		}
		return &abi.MintProposal{
			Id:         new(big.Int).SetUint64(proposalID),
			Name:       "TestAgent",
			Role:       "expert",
			Owner:      common_tests.RandomEthAddress(),
			Creator:    common_tests.RandomEthAddress(),
			Addresses:  addresses,
			Images:     []abi.KeyValue{{Key: "avatar", Value: "ipfs://test"}},
			Attributes: []abi.KeyValue{{Key: "bio", Value: "test bio"}},
		}
	}

	newEvent := func(collectionAddress common.Address, proposalID uint64, blockNumber uint64) *abi.AgentCollectionV1MintProposalCreated {
		return &abi.AgentCollectionV1MintProposalCreated{
			From:       common_tests.RandomEthAddress(),
			ProposalId: new(big.Int).SetUint64(proposalID),
			Raw: types.Log{
				Address:     collectionAddress,
				BlockNumber: blockNumber,
			},
		}
	}

	t.Run("Given retrieving mint proposal fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 1, 100)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(100)).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped requester error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get mint proposal")
		})
	})

	t.Run("Given proposal carries only unknown address types", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 2, 200)
			proposal := newMintProposal(2, []common.Address{})
			proposal.Addresses = []abi.AgentAddress{{AddressType: "solana", Value: "unknown-type-value"}}
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(200)).
				Return(proposal, nil)

			return event
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "no known agent address found")
		})
	})

	t.Run("Given agent lookup fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 3, 300)
			proposal := newMintProposal(3, []common.Address{agentWallet})
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(300)).
				Return(proposal, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(agentWallet.Hex())).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped repository error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent")
		})
	})

	t.Run("Given agent already exists and is already minted", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 8, 800)
			proposal := newMintProposal(8, []common.Address{agentWallet})
			tokenId := uint64(42)
			mintedAgent := model_builders.NewAgentBuilder().
				WithCollectionAddress(collectionAddress.Hex()).
				WithCollectionTokenId(&tokenId).
				WithProposalId(nil).
				Build()

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(800)).
				Return(proposal, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(agentWallet.Hex())).
				Return(model_builders.NewAgentAddressBuilder().WithAgentId(mintedAgent.ID).Build(), nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetById(mintedAgent.ID).
				Return(mintedAgent, nil)

			return event
		}

		t.Run("Should treat the event as a no-op and not error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given lookup of already-minted agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 9, 900)
			proposal := newMintProposal(9, []common.Address{agentWallet})

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(900)).
				Return(proposal, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(agentWallet.Hex())).
				Return(model_builders.NewAgentAddressBuilder().Build(), nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped repository error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent by address")
		})
	})

	t.Run("Given agent does not exist and is external", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) (*abi.AgentCollectionV1MintProposalCreated, common.Address) {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 5, 500)
			proposal := newMintProposal(5, []common.Address{agentWallet})
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(500)).
				Return(proposal, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(agentWallet.Hex())).
				Return(nil, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeBtc, "bc1qexampleonly").
				Return(nil, nil)
			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(collectionAddress.Hex()).
				Return(model_builders.NewAgentCollectionBuilder().WithName("Genesis").Build(), nil)
			suite.decentralizedDomainNameBuilderMock.
				EXPECT().
				BuildDomainName("TestAgent", "Genesis").
				Return("testagent.genesis.eth")
			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					assert.Equal(t, collectionAddress.Hex(), saved.CollectionAddress)
					assert.Equal(t, "TestAgent", saved.Name)
					assert.Equal(t, "testagent.genesis.eth", saved.EnsDomain)
					assert.NotNil(t, saved.Owner)
					assert.Equal(t, map[string]string{"avatar": "ipfs://test"}, saved.Images)
					assert.Equal(t, map[string]string{"bio": "test bio"}, saved.Attributes)
					if assert.NotNil(t, saved.ProposalId) {
						assert.Equal(t, uint64(5), *saved.ProposalId)
					}
					return saved, nil
				})
			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(gomock.Any(), map[enums.AddressType]string{
					enums.AddressTypeEvm: strings.ToLower(agentWallet.Hex()),
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil)

			return event, agentWallet
		}

		t.Run("Should create external agent with proposal id and save", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event, _ := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})

	t.Run("Given saving agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(collectionAddress.Hex()).
				WithCollectionTokenId(nil).
				Build()
			event := newEvent(collectionAddress, 6, 600)
			proposal := newMintProposal(6, []common.Address{agentWallet})
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(600)).
				Return(proposal, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(agentWallet.Hex())).
				Return(model_builders.NewAgentAddressBuilder().WithAgentId(agent.ID).Build(), nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetById(agent.ID).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to save agent")
		})
	})

	t.Run("Given existing agent and event is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(collectionAddress.Hex()).
				WithCollectionTokenId(nil).
				Build()
			event := newEvent(collectionAddress, 7, 700)
			proposal := newMintProposal(7, []common.Address{agentWallet})
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetMintProposal(collectionAddress, gomock.Any(), uint64(700)).
				Return(proposal, nil)
			suite.agentAddressRepositoryMock.
				EXPECT().
				GetByTypeAndAddress(enums.AddressTypeEvm, strings.ToLower(agentWallet.Hex())).
				Return(model_builders.NewAgentAddressBuilder().WithAgentId(agent.ID).Build(), nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetById(agent.ID).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					if assert.NotNil(t, saved.ProposalId) {
						assert.Equal(t, uint64(7), *saved.ProposalId)
					}
					return saved, nil
				})
			suite.agentAddressRepositoryMock.
				EXPECT().
				ReplaceAllForAgent(agent.ID, map[enums.AddressType]string{
					enums.AddressTypeEvm: strings.ToLower(agentWallet.Hex()),
					enums.AddressTypeBtc: "bc1qexampleonly",
				}).
				Return(nil)
			return event
		}

		t.Run("Should set proposal id and save", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.NoError(t, err)
		})
	})
}
