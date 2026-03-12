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
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingMintProposalCreatedEventTestingSuite struct {
	agentRepositoryMock          *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRequesterMock *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.MintProposalCreatedEventHandler
}

func WhenHandlingMintProposalCreatedEventBeforeEach(t *testing.T) *WhenHandlingMintProposalCreatedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)

	sut := services.NewMintProposalCreatedEventHandler(
		zap.NewNop(),
		agentRepositoryMock,
		agentCollectionRequesterMock,
	)

	return &WhenHandlingMintProposalCreatedEventTestingSuite{
		agentRepositoryMock:          agentRepositoryMock,
		agentCollectionRequesterMock: agentCollectionRequesterMock,

		sut: sut,
	}
}

func TestWhenHandlingMintProposalCreatedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionAddress common.Address, proposalID uint64) *abi.AgentCollectionV1MintProposalCreated {
		return &abi.AgentCollectionV1MintProposalCreated{
			From:       common_tests.RandomEthAddress(),
			Signer:     common_tests.RandomEthAddress(),
			ProposalId: new(big.Int).SetUint64(proposalID),
			Raw: types.Log{
				Address: collectionAddress,
			},
		}
	}

	t.Run("Given retrieving proposal wallets fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 1)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(collectionAddress, gomock.Any()).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped requester error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent wallets")
		})
	})

	t.Run("Given proposal has no wallets", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 2)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(collectionAddress, gomock.Any()).
				Return([]common.Address{}, nil)

			return event
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "no agent wallets found")
		})
	})

	t.Run("Given agent lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 3)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(collectionAddress, gomock.Any()).
				Return([]common.Address{agentWallet}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(collectionAddress.Hex(), agentWallet.Hex()).
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

	t.Run("Given no agent matches proposal wallet", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 4)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(collectionAddress, gomock.Any()).
				Return([]common.Address{agentWallet}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(collectionAddress.Hex(), agentWallet.Hex()).
				Return(nil, nil)

			return event
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintProposalCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "no agent found")
		})
	})

	t.Run("Given saving agent fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().Build()
			event := newEvent(collectionAddress, 5)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(collectionAddress, gomock.Any()).
				Return([]common.Address{agentWallet}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(collectionAddress.Hex(), agentWallet.Hex()).
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

	t.Run("Given event is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintProposalCreatedEventTestingSuite) *abi.AgentCollectionV1MintProposalCreated {
			collectionAddress := common_tests.RandomEthAddress()
			agentWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().Build()
			event := newEvent(collectionAddress, 6)
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(collectionAddress, gomock.Any()).
				Return([]common.Address{agentWallet}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(collectionAddress.Hex(), agentWallet.Hex()).
				Return(agent, nil)
			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				DoAndReturn(func(saved *models.Agent) (*models.Agent, error) {
					if assert.NotNil(t, saved.ProposalId) {
						assert.Equal(t, uint64(6), *saved.ProposalId)
					}
					return saved, nil
				})
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
