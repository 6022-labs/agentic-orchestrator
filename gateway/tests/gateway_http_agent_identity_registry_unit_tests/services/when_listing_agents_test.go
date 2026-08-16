package services_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_identity_registry_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentsTestingSuite struct {
	sut *services.AgentRegistryRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenListingAgentsBeforeEach(t *testing.T) *WhenListingAgentsTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenListingAgentsTestingSuite{
		sut:        services.NewAgentRegistryRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenListingAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given the registry returns agents", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) uuid.UUID {
			agentId := uuid.New()

			suite.mockClient.EXPECT().GetAgentsWithResponse(gomock.Any(), &clients.GetAgentsParams{
				Page:            lo.ToPtr(2),
				ElementsPerPage: lo.ToPtr(10),
			}).Return(&clients.GetAgentsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &clients.ResponsesListAgentsResponse{
					Page:            lo.ToPtr(2),
					ElementsPerPage: lo.ToPtr(10),
					TotalElements:   lo.ToPtr(11),
					Agents: &[]clients.ResponsesAgentResponse{{
						Id:                lo.ToPtr(agentId.String()),
						Name:              lo.ToPtr("Ada"),
						Attributes:        lo.ToPtr(map[string]string{"role": "expert"}),
						EnsDomain:         lo.ToPtr("ada.experts.80002.6022.eth"),
						CollectionAddress: lo.ToPtr("0xabc"),
						CollectionTokenId: lo.ToPtr(7),
					}},
				},
			}, nil).Times(1)

			return agentId
		}

		t.Run("Should map the generated response to the use case response", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			agentId := initSuite(suite)

			result, err := suite.sut.ListAgents(&interfaces.ListRegistryAgentsRequest{
				Page:            lo.ToPtr(uint(2)),
				ElementsPerPage: lo.ToPtr(uint(10)),
			})

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, uint(2), result.Page)
				assert.Equal(t, uint(10), result.ElementsPerPage)
				assert.Equal(t, uint(11), result.TotalElements)
				if assert.Len(t, result.Agents, 1) {
					agent := result.Agents[0]
					assert.Equal(t, agentId, agent.ID)
					assert.Equal(t, "Ada", agent.Name)
					assert.Equal(t, "expert", agent.Attributes["role"])
					assert.Equal(t, "ada.experts.80002.6022.eth", agent.EnsDomain)
					assert.Equal(t, "0xabc", agent.CollectionAddress)
					if assert.NotNil(t, agent.CollectionTokenId) {
						assert.Equal(t, uint64(7), *agent.CollectionTokenId)
					}
				}
			}
		})
	})

	t.Run("Given the registry returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) {
			suite.mockClient.EXPECT().GetAgentsWithResponse(gomock.Any(), gomock.Any()).Return(&clients.GetAgentsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
				Body:         []byte("boom"),
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ListAgents(&interfaces.ListRegistryAgentsRequest{})

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) {
			suite.mockClient.EXPECT().GetAgentsWithResponse(gomock.Any(), gomock.Any()).Return(nil, errors.New("registry down")).Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ListAgents(&interfaces.ListRegistryAgentsRequest{})

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
