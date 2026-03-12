package services_test

import (
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
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenHandlingMintedEventTestingSuite struct {
	orchestratorSettings               *settings.OrchestratorSettings
	agentRepositoryMock                *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRepositoryMock      *repositories_mocks.MockAgentCollectionRepositoryInterface
	decentralizedDomainNameBuilder     *interfaces_mocks.MockDecentralizedDomainNameBuilderInterface
	agentCollectionBlockchainRequester *interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface

	sut *services.MintedEventHandler
}

func WhenHandlingMintedEventBeforeEach(t *testing.T) *WhenHandlingMintedEventTestingSuite {
	mockController := gomock.NewController(t)

	orchestratorPublicKey := common_tests.RandomEthAddress()
	orchestratorSettings := &settings.OrchestratorSettings{PublicKey: orchestratorPublicKey}

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	decentralizedDomainNameBuilder := interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(mockController)
	agentCollectionBlockchainRequester := interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)

	sut := services.NewMintedEventHandler(
		zap.NewNop(),
		orchestratorSettings,
		agentRepositoryMock,
		agentCollectionRepositoryMock,
		decentralizedDomainNameBuilder,
		agentCollectionBlockchainRequester,
	)

	return &WhenHandlingMintedEventTestingSuite{
		orchestratorSettings:               orchestratorSettings,
		agentRepositoryMock:                agentRepositoryMock,
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
			Signer:  common_tests.RandomEthAddress(),
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

	t.Run("Given agent wallet lookup fails with unexpected error", func(t *testing.T) {
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
					Signer:  common_tests.RandomEthAddress(),
					Creator: common_tests.RandomEthAddress(),
					Wallets: []common.Address{wallet},
				}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetByWalletAddress(wallet.Hex()).
				Return(nil, assert.AnError)
			return event
		}

		t.Run("Should return wrapped wallet lookup error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to find agent by wallet addresses")
		})
	})

	t.Run("Given internal signer agent does not exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMintedEventTestingSuite) *abi.AgentCollectionV1Minted {
			collectionAddress := common_tests.RandomEthAddress()
			wallet := common_tests.RandomEthAddress()
			event := newEvent(collectionAddress, 3)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentInformation(collectionAddress, gomock.Any()).
				Return(&abi.AgentInformation{
					Name:    "Agent",
					Role:    "expert",
					Signer:  suite.orchestratorSettings.PublicKey,
					Creator: common_tests.RandomEthAddress(),
					Wallets: []common.Address{wallet},
				}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetByWalletAddress(wallet.Hex()).
				Return(nil, gorm.ErrRecordNotFound)
			return event
		}

		t.Run("Should return explicit internal agent error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMintedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "must be created by this orchestrator")
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
					Signer:  common_tests.RandomEthAddress(),
					Creator: common_tests.RandomEthAddress(),
					Wallets: []common.Address{wallet},
				}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetByWalletAddress(wallet.Hex()).
				Return(nil, gorm.ErrRecordNotFound)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(collectionAddress, gomock.Any()).
				Return(common_tests.RandomEthAddress(), nil)
			suite.agentCollectionRepositoryMock.
				EXPECT().
				GetByAddress(collectionAddress.Hex()).
				Return(nil, gorm.ErrRecordNotFound)

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
					Name:       "Alice",
					Role:       "human",
					Signer:     common_tests.RandomEthAddress(),
					Creator:    common_tests.RandomEthAddress(),
					Wallets:    []common.Address{wallet},
					Images:     []abi.KeyValue{{Key: "avatar", Value: "ipfs://avatar"}},
					Attributes: []abi.KeyValue{{Key: "skill", Value: "analysis"}},
				}, nil)
			suite.agentRepositoryMock.
				EXPECT().
				GetByWalletAddress(wallet.Hex()).
				Return(agent, nil)
			suite.agentCollectionBlockchainRequester.
				EXPECT().
				GetAgentOwner(collectionAddress, gomock.Any()).
				Return(owner, nil)
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
					assert.Equal(t, wallet.Hex(), saved.WalletAddress)
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
