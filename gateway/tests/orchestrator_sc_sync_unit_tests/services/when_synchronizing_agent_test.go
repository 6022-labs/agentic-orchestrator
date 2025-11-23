package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSynchronizingAgentTestingSuite struct {
	sut *services.AgentSynchronizer

	orchestratorSettings               *settings.OrchestratorSettings
	agentRepositoryMock                *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRequesterMock       *interfaces_mocks.MockAgentCollectionRequesterInterface
	agentCollectionRepositoryMock      *repositories_mocks.MockAgentCollectionRepositoryInterface
	decentralizedDomainNameBuilderMock *interfaces_mocks.MockDecentralizedDomainNameBuilderInterface
}

func WhenSynchronizingAgentBeforeEach(t *testing.T) *WhenSynchronizingAgentTestingSuite {
	mockController := gomock.NewController(t)

	orchestratorPublicKey := common_tests.RandomEthAddress()
	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: orchestratorPublicKey,
	}

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	decentralizedDomainNameBuilderMock := interfaces_mocks.NewMockDecentralizedDomainNameBuilderInterface(mockController)

	sut := services.NewAgentSynchronizer(
		orchestratorSettings,
		agentRepositoryMock,
		agentCollectionRequesterMock,
		agentCollectionRepositoryMock,
		decentralizedDomainNameBuilderMock,
	)

	return &WhenSynchronizingAgentTestingSuite{
		sut: sut,

		agentRepositoryMock:                agentRepositoryMock,
		orchestratorSettings:               orchestratorSettings,
		agentCollectionRequesterMock:       agentCollectionRequesterMock,
		agentCollectionRepositoryMock:      agentCollectionRepositoryMock,
		decentralizedDomainNameBuilderMock: decentralizedDomainNameBuilderMock,
	}
}

func TestWhenSynchronizingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting agent information on-chain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingAgentTestingSuite) {
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			suite := WhenSynchronizingAgentBeforeEach(t)
			initSuite(suite)

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			err := suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while finding agent by wallet addresses", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingAgentTestingSuite) {
			agentInformation := abi.AgentInformation{
				Role:    enums.AgentRoleExpert.String(),
				Wallets: []common.Address{common.HexToAddress("0x1D1BF0f293068F3278a8a076d66a1883B4980864")},
			}
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), gomock.Any()).
				Return(&agentInformation, nil)

			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			suite := WhenSynchronizingAgentBeforeEach(t)
			initSuite(suite)

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			err := suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not found but orchestrator must be the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingAgentTestingSuite) {
			agentInformation := abi.AgentInformation{
				Role:    enums.AgentRoleExpert.String(),
				Signer:  suite.orchestratorSettings.PublicKey,
				Wallets: []common.Address{common.HexToAddress("0x1D1BF0f293068F3278a8a076d66a1883B4980864")},
			}
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), gomock.Any()).
				Return(&agentInformation, nil)

			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(gomock.Any()).Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			suite := WhenSynchronizingAgentBeforeEach(t)

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			initSuite(suite)

			err := suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not found and orchestrator is not the signer", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingAgentTestingSuite) *abi.AgentInformation {
			owner := common.HexToAddress("0x2332cc64591301E10D08f503A75B05D0903D6fe2")
			agentInformation := abi.AgentInformation{
				Name: "agent-test",
				Role: enums.AgentRoleFacilitator.String(),
				Images: []abi.KeyValue{
					{Key: "default", Value: "bafybeia25xzo4k4nipuw6cyml4f4x3b3h6lneyvch6a4b3rl2q7e6hpgra"},
				},
				Attributes: []abi.KeyValue{
					{Key: "mission", Value: "test-mission"},
				},
				Creator: owner,
				Wallets: []common.Address{common.HexToAddress("0x1D1BF0f293068F3278a8a076d66a1883B4980864")},
			}
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), gomock.Any()).
				Return(&agentInformation, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentOwner(gomock.Any(), gomock.Any()).
				Return(&owner, nil)

			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(gomock.Any()).Return(nil, nil)

			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(gomock.Any()).Return(&models.AgentCollection{Name: "protocol6022"}, nil)

			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(agentInformation.Name, "protocol6022").Return(agentInformation.Name + ".protocol6022.eth")

			return &agentInformation
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentBeforeEach(t)

			suite.agentRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, nil)

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			initSuite(suite)

			err := suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
			assert.NoError(t, err)
		})

		t.Run("Should create external agent", func(t *testing.T) {
			t.Parallel()

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			suite := WhenSynchronizingAgentBeforeEach(t)
			agentInformation := initSuite(suite)
			suite.agentRepositoryMock.EXPECT().Save(gomock.Any()).DoAndReturn(func(agent *models.Agent) (*models.Agent, error) {
				assert.False(t, agent.IsInternal)
				assert.Equal(t, agentInformation.Name, agent.Name)
				assert.Equal(t, agentInformation.Role, agent.Role.String())
				assert.Equal(t, agentInformation.Wallets[0].Hex(), agent.WalletAddress)
				assert.Equal(t, agentInformation.Creator.Hex(), agent.Creator)
				assert.Equal(t, agentInformation.Creator.Hex(), *agent.Owner)
				assert.Equal(t, agentCollectionTokenId, *agent.CollectionTokenId)
				assert.Equal(t, agentCollectionAddress, agent.CollectionAddress)
				assert.Equal(t, agentInformation.Name+".protocol6022.eth", agent.EnsDomain)
				assert.Equal(t, 1, len(agent.Attributes))
				assert.Equal(t, agentInformation.Attributes[0].Value, agent.Attributes["mission"])
				assert.Equal(t, 1, len(agent.Images))
				assert.Equal(t, "bafybeia25xzo4k4nipuw6cyml4f4x3b3h6lneyvch6a4b3rl2q7e6hpgra", agent.Images["default"])
				return agent, nil
			})

			suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
		})
	})

	t.Run("Given creating agent is found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingAgentTestingSuite) *abi.AgentInformation {
			agentInformation := abi.AgentInformation{
				Name: "agent-test",
				Role: enums.AgentRoleFacilitator.String(),
				Images: []abi.KeyValue{
					{Key: "default", Value: "bafybeia25xzo4k4nipuw6cyml4f4x3b3h6lneyvch6a4b3rl2q7e6hpgra"},
				},
				Attributes: []abi.KeyValue{
					{Key: "mission", Value: "test-mission"},
				},
				Creator: common.HexToAddress("0x2332cc64591301E10D08f503A75B05D0903D6fe2"),
				Wallets: []common.Address{common.HexToAddress("0x1D1BF0f293068F3278a8a076d66a1883B4980864")},
			}
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentInformation(gomock.Any(), gomock.Any()).
				Return(&agentInformation, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentOwner(gomock.Any(), gomock.Any()).
				Return(&agentInformation.Creator, nil)

			agent := model_builders.
				NewAgentBuilder().
				WithCollectionTokenId(nil).
				WithWalletAddress(agentInformation.Wallets[0].Hex()).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByWalletAddress(gomock.Any()).Return(agent, nil)

			suite.agentCollectionRepositoryMock.EXPECT().GetByAddress(gomock.Any()).Return(&models.AgentCollection{Name: "protocol6022"}, nil)

			suite.decentralizedDomainNameBuilderMock.EXPECT().BuildDomainName(agentInformation.Name, "protocol6022").Return(agentInformation.Name + ".protocol6022.eth")

			return &agentInformation
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentBeforeEach(t)
			initSuite(suite)
			suite.agentRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, nil)

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			err := suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
			assert.NoError(t, err)
		})

		t.Run("Should update existing agent", func(t *testing.T) {
			t.Parallel()

			agentCollectionAddress := common_tests.RandomEthAddress().Hex()
			agentCollectionTokenId := uint64(1)

			suite := WhenSynchronizingAgentBeforeEach(t)
			agentInformation := initSuite(suite)
			suite.agentRepositoryMock.EXPECT().Save(gomock.Any()).DoAndReturn(func(agent *models.Agent) (*models.Agent, error) {
				assert.Equal(t, agentInformation.Name, agent.Name)
				assert.Equal(t, agentInformation.Role, agent.Role.String())
				assert.Equal(t, agentInformation.Creator.Hex(), agent.Creator)
				assert.Equal(t, agentCollectionTokenId, *agent.CollectionTokenId)
				assert.Equal(t, agentCollectionAddress, agent.CollectionAddress)
				assert.Equal(t, agentInformation.Name+".protocol6022.eth", agent.EnsDomain)
				assert.Equal(t, 1, len(agent.Attributes))
				assert.Equal(t, agentInformation.Attributes[0].Value, agent.Attributes["mission"])
				assert.Equal(t, 1, len(agent.Images))
				assert.Equal(t, "bafybeia25xzo4k4nipuw6cyml4f4x3b3h6lneyvch6a4b3rl2q7e6hpgra", agent.Images["default"])
				return agent, nil
			})

			err := suite.sut.Sync(agentCollectionAddress, agentCollectionTokenId)
			assert.NoError(t, err)
		})
	})
}
