package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenPairingProposalToAgentTestingSuite struct {
	sut *use_cases.PairProposalToAgent

	agentRepositoryMock          *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRequesterMock *interfaces_mocks.MockAgentCollectionRequesterInterface
}

func WhenPairingProposalToAgentBeforeEach(t *testing.T) *WhenPairingProposalToAgentTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)

	sut := use_cases.NewPairProposalToAgent(
		agentRepositoryMock,
		agentCollectionRequesterMock,
	)

	return &WhenPairingProposalToAgentTestingSuite{
		sut: sut,

		agentRepositoryMock:          agentRepositoryMock,
		agentCollectionRequesterMock: agentCollectionRequesterMock,
	}
}

func generateValidPairProposalToAgentRequest() *requests.PairProposalToAgentRequest {
	return &requests.PairProposalToAgentRequest{
		CollectionAddress: "0x1234567890abcdef1234567890abcdef12345678",
		ProposalId:        1,
	}
}

func TestWhenPairingProposalToAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collection address is not a valid Ethereum address", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			request := generateValidPairProposalToAgentRequest()
			request.CollectionAddress = "invalid-address"

			err := suite.sut.Pair(request)

			assert.Error(t, err)
		})

		t.Run("Because proposal ID is zero", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			request := generateValidPairProposalToAgentRequest()
			request.ProposalId = 0

			err := suite.sut.Pair(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent wallets", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingProposalToAgentTestingSuite) {
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			initSuite(suite)

			request := generateValidPairProposalToAgentRequest()

			err := suite.sut.Pair(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no agent wallets found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingProposalToAgentTestingSuite) {
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(gomock.Any(), gomock.Any()).
				Return([]common.Address{}, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			initSuite(suite)

			request := generateValidPairProposalToAgentRequest()

			err := suite.sut.Pair(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent by collection address and wallet address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingProposalToAgentTestingSuite) {
			walletAddress := common_tests.RandomEthAddress()
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(gomock.Any(), gomock.Any()).
				Return([]common.Address{walletAddress}, nil)

			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			initSuite(suite)

			request := generateValidPairProposalToAgentRequest()

			err := suite.sut.Pair(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no agent found for collection address and wallet address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingProposalToAgentTestingSuite) {
			walletAddress := common_tests.RandomEthAddress()
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(gomock.Any(), gomock.Any()).
				Return([]common.Address{walletAddress}, nil)

			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(gomock.Any(), gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			initSuite(suite)

			request := generateValidPairProposalToAgentRequest()

			err := suite.sut.Pair(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent found and successfully paired", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingProposalToAgentTestingSuite) {
			walletAddress := common_tests.RandomEthAddress()

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentWalletsOfProposal(gomock.Any(), gomock.Any()).
				Return([]common.Address{walletAddress}, nil)

			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndWalletAddress(gomock.Any(), gomock.Any()).
				Return(agent, nil)

			suite.agentRepositoryMock.
				EXPECT().
				Save(agent).
				Return(agent, nil)
		}

		t.Run("Should successfully pair proposal to agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingProposalToAgentBeforeEach(t)
			initSuite(suite)

			request := generateValidPairProposalToAgentRequest()

			err := suite.sut.Pair(request)

			assert.NoError(t, err)
		})
	})
}
