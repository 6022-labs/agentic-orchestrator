package use_cases_test

import (
	"context"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.SaveAgentPromptTemplate

	mockAgentRegistryRequester        *interfaces_mocks.MockAgentRegistryRequesterInterface
	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(
		mockController,
	)

	sut := use_cases.NewSaveAgentPromptTemplate(mockAgentRegistryRequester, mockAgentPromptTemplateRepository)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentRegistryRequester:        mockAgentRegistryRequester,
		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
	}
}

func generateValidSaveAgentPromptTemplateRequest() *requests.SaveAgentPromptTemplateRequest {
	return &requests.SaveAgentPromptTemplateRequest{
		GetAgentPromptTemplateRequest: requests.GetAgentPromptTemplateRequest{
			CollectionAddress: "0xabc00000000000000000000000000000000000de",
			CollectionTokenId: 1,
			Key:               "clone-system",
		},
		Content: "You are {{agentName}}.",
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress has no 0x prefix", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.CollectionAddress = "abc00000000000000000000000000000000000de"

			_, err := suite.sut.Save(context.Background(), request)

			var validationErr *cb_errors.ValidationError
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "collectionAddress", validationErr.Field)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0xabc00000000000000000000000000000000000de", uint64(1)).
				Return(model_builders.NewRegistryAgentBuilder().
					WithRole("clone").
					WithCollection("0xabc00000000000000000000000000000000000de", 1).
					Build(), nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), "0xabc00000000000000000000000000000000000de", uint64(1), gomock.Any()).
				Return(nil, nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, entity *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error) {
					return entity, nil
				}).
				Times(1)
		}

		t.Run("Should fetch, look up and persist the lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.CollectionAddress = "0xAbC00000000000000000000000000000000000dE"

			saved, err := suite.sut.Save(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xabc00000000000000000000000000000000000de", saved.CollectionAddress)
		})
	})

	t.Run("Given registry lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request)

			_, err := suite.sut.Save(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not registered", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should reject with a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request)

			_, err := suite.sut.Save(context.Background(), request)

			var notFoundErr *cb_errors.NotFoundError
			require.ErrorAs(t, err, &notFoundErr)
			assert.Equal(t, "agent_not_found", notFoundErr.Code)
		})
	})

	t.Run("Given agent role does not allow the key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest, role string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithRole(role).
						WithCollection(request.CollectionAddress, request.CollectionTokenId).
						Build(),
					nil,
				).
				Times(1)
		}

		t.Run("Should reject clone keys for a facilitator agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request, "facilitator")

			_, err := suite.sut.Save(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})

		t.Run("Should reject clone keys for a human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request, "human")

			_, err := suite.sut.Save(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})

		t.Run("Should reject facilitator keys for a human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Key = "facilitator-system"
			initSuite(suite, request, "human")

			_, err := suite.sut.Save(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})

		t.Run("Should reject non-expert keys for an unknown role", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request, "banana")

			_, err := suite.sut.Save(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given agent role falls back to expert", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest, role string) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithRole(role).
						WithCollection(request.CollectionAddress, request.CollectionTokenId).
						Build(),
					nil,
				).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), request.CollectionAddress, request.CollectionTokenId, request.Key).
				Return(nil, nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, entity *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error) {
					return entity, nil
				}).
				Times(1)
		}

		t.Run("Should allow expert keys for an unknown role", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Key = "expert-system"
			initSuite(suite, request, "banana")

			_, err := suite.sut.Save(context.Background(), request)

			assert.NoError(t, err)
		})

		t.Run("Should allow expert keys for a human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Key = "expert-default"
			initSuite(suite, request, "human")

			_, err := suite.sut.Save(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given saving template fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithRole("clone").
						WithCollection(request.CollectionAddress, request.CollectionTokenId).
						Build(),
					nil,
				).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request)

			_, err := suite.sut.Save(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given template does not exist yet", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithRole("clone").
						WithCollection(request.CollectionAddress, request.CollectionTokenId).
						Build(),
					nil,
				).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), request.CollectionAddress, request.CollectionTokenId, request.Key).
				Return(nil, nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, entity *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error) {
					assert.Equal(t, request.CollectionAddress, entity.CollectionAddress)
					assert.Equal(t, request.CollectionTokenId, entity.CollectionTokenId)
					assert.Equal(t, request.Key, entity.Key)
					assert.Equal(t, request.Content, entity.Content)
					return entity, nil
				}).
				Times(1)
		}

		t.Run("Should create the template", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request)

			saved, err := suite.sut.Save(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, saved.Content)
		})
	})

	t.Run("Given template already exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) *models.AgentPromptTemplate {
			existing := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				WithKey(request.Key).
				WithContent("old content").
				Build()

			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithRole("clone").
						WithCollection(request.CollectionAddress, request.CollectionTokenId).
						Build(),
					nil,
				).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), request.CollectionAddress, request.CollectionTokenId, request.Key).
				Return(existing, nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any(), existing).
				Return(existing, nil).
				Times(1)

			return existing
		}

		t.Run("Should overwrite the existing template content", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			existing := initSuite(suite, request)

			saved, err := suite.sut.Save(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, existing.ID, saved.ID)
			assert.Equal(t, request.Content, saved.Content)
		})
	})
}
