package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingMintedAgentByEnsDomainTestingSuite struct {
	sut *use_cases.GetMintedAgentByEnsDomain

	mockAgentRepository        *repositories_mocks.MockAgentRepositoryInterface
	mockAgentAddressRepository *repositories_mocks.MockAgentAddressRepositoryInterface
}

func WhenGettingMintedAgentByEnsDomainBeforeEach(t *testing.T) *WhenGettingMintedAgentByEnsDomainTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockAgentAddressRepository := repositories_mocks.NewMockAgentAddressRepositoryInterface(mockController)

	sut := use_cases.NewGetMintedAgentByEnsDomain(mockAgentRepository, mockAgentAddressRepository)

	return &WhenGettingMintedAgentByEnsDomainTestingSuite{
		sut:                        sut,
		mockAgentRepository:        mockAgentRepository,
		mockAgentAddressRepository: mockAgentAddressRepository,
	}
}

func TestWhenGettingMintedAgentByEnsDomain(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingMintedAgentByEnsDomainTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				GetMintedByEnsDomain(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingMintedAgentByEnsDomainBeforeEach(t)
			initSuite(basics)

			req := &requests.GetAgentByEnsDomainRequest{
				EnsDomain: "alice.6022.eth",
			}
			resp, err := basics.sut.Get(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingMintedAgentByEnsDomainTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				GetMintedByEnsDomain(gomock.Any()).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return nil response and nil error", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingMintedAgentByEnsDomainBeforeEach(t)
			initSuite(basics)

			req := &requests.GetAgentByEnsDomainRequest{
				EnsDomain: "unknown.6022.eth",
			}
			resp, err := basics.sut.Get(req)

			assert.Nil(t, resp)
			assert.Nil(t, err)
		})
	})

	t.Run("Given agent is fetched successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the agent response", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingMintedAgentByEnsDomainBeforeEach(t)

			agent := model_builders.NewAgentBuilder().
				WithEnsDomain("alice.6022.eth").
				Build()

			basics.mockAgentRepository.EXPECT().
				GetMintedByEnsDomain(agent.EnsDomain).
				Return(agent, nil).
				Times(1)

			basics.mockAgentAddressRepository.EXPECT().
				GetAllByAgentId(agent.ID).
				Return([]*models.AgentAddress{
					model_builders.NewAgentAddressBuilder().WithAgentId(agent.ID).WithAddress("0x2f0ddd93040c19863527b0553feceb696862ea7d").Build(),
				}, nil).
				Times(1)

			req := &requests.GetAgentByEnsDomainRequest{
				EnsDomain: agent.EnsDomain,
			}
			resp, err := basics.sut.Get(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, agent.ID, resp.ID)
			assert.Equal(t, agent.Name, resp.Name)
			assert.Equal(t, agent.Role, resp.Role)
			assert.Equal(t, agent.Owner, resp.Owner)
			assert.Equal(t, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
			}, resp.Addresses)
			assert.Equal(t, agent.CollectionAddress, resp.CollectionAddress)
			assert.Equal(t, agent.CollectionTokenId, resp.CollectionTokenId)
			assert.Equal(t, agent.EnsDomain, resp.EnsDomain)
		})
	})
}
