package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.GetAgentPromptTemplate

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)

	sut := use_cases.NewGetAgentPromptTemplate(mockAgentPromptTemplateRepository)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
	}
}

func generateValidGetAgentPromptTemplateRequest() *requests.GetAgentPromptTemplateRequest {
	return &requests.GetAgentPromptTemplateRequest{
		CollectionAddress: "0xabc",
		CollectionTokenId: "1",
		Key:               "clone-system",
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.CollectionAddress = ""

			_, err := suite.sut.Get(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})

		t.Run("Because collectionTokenId is not a number", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.CollectionTokenId = "abc"

			_, err := suite.sut.Get(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})

		t.Run("Because key is unknown", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.Key = "unknown-key"

			_, err := suite.sut.Get(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given template does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(generateValidGetAgentPromptTemplateRequest())

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given template exists", func(t *testing.T) {
		t.Parallel()

		template := model_builders.NewAgentPromptTemplateBuilder().
			WithCollectionAddress("0xabc").
			WithCollectionTokenId(1).
			Build()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey("0xabc", uint64(1), template.Key).
				Return(template, nil).
				Times(1)
		}

		t.Run("Should return the template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Get(generateValidGetAgentPromptTemplateRequest())

			assert.NoError(t, err)
			assert.Equal(t, template.ID, got.ID)
		})
	})

	t.Run("Given repository fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(generateValidGetAgentPromptTemplateRequest())

			assert.Error(t, err)
		})
	})
}
