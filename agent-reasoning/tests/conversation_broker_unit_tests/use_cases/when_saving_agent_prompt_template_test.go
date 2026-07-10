package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.SaveAgentPromptTemplate

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)

	sut := use_cases.NewSaveAgentPromptTemplate(mockAgentPromptTemplateRepository)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
	}
}

func generateValidSaveAgentPromptTemplateRequest() *requests.SaveAgentPromptTemplateRequest {
	return &requests.SaveAgentPromptTemplateRequest{
		GetAgentPromptTemplateRequest: requests.GetAgentPromptTemplateRequest{
			CollectionAddress: "0xabc",
			CollectionTokenId: "1",
			Key:               "clone-system",
		},
		Content: "You are {{agentName}}.",
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because content is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Content = "  "

			_, err := suite.sut.Save(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given template does not exist yet", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(request.CollectionAddress, request.ParsedCollectionTokenId(), request.ParsedKey()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(entity *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error) {
					assert.Equal(t, request.CollectionAddress, entity.CollectionAddress)
					assert.Equal(t, request.ParsedCollectionTokenId(), entity.CollectionTokenId)
					assert.Equal(t, request.ParsedKey(), entity.Key)
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

			saved, err := suite.sut.Save(request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, saved.Content)
		})
	})

	t.Run("Given template already exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) *models.AgentPromptTemplate {
			existing := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.ParsedCollectionTokenId()).
				WithKey(request.ParsedKey()).
				WithContent("old content").
				Build()

			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(request.CollectionAddress, request.ParsedCollectionTokenId(), request.ParsedKey()).
				Return(existing, nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(existing).
				Return(existing, nil).
				Times(1)

			return existing
		}

		t.Run("Should overwrite the existing template content", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			existing := initSuite(suite, request)

			saved, err := suite.sut.Save(request)

			assert.NoError(t, err)
			assert.Equal(t, existing.ID, saved.ID)
			assert.Equal(t, request.Content, saved.Content)
		})
	})

	t.Run("Given saving template fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Save(generateValidSaveAgentPromptTemplateRequest())

			assert.Error(t, err)
		})
	})
}
