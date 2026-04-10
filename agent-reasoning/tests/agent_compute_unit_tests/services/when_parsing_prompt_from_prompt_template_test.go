package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenParsingPromptFromPromptTemplateTestingSuite struct {
	sut                          *services.PromptFromPromptTemplateParser
	agentPromptTemplateRequester *interfaces_mocks.MockAgentPromptTemplateRequesterInterface
}

func WhenParsingPromptFromPromptTemplateBeforeEach(t *testing.T) *WhenParsingPromptFromPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	agentPromptTemplateRequesterMock := interfaces_mocks.NewMockAgentPromptTemplateRequesterInterface(mockController)

	sut := services.NewPromptFromPromptTemplateParser(agentPromptTemplateRequesterMock)

	return &WhenParsingPromptFromPromptTemplateTestingSuite{
		sut:                          sut,
		agentPromptTemplateRequester: agentPromptTemplateRequesterMock,
	}
}

func TestWhenParsingPromptFromPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("ParseExpertSystem", func(t *testing.T) {
		t.Parallel()

		t.Run("Given template requester returns error", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingPromptFromPromptTemplateBeforeEach(t)

			agentId := uuid.New()
			parameters := map[string]any{"Name": "test"}

			suite.agentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertSystem).
				Return(nil, assert.AnError)

			result, err := suite.sut.ParseExpertSystem(agentId, parameters)

			assert.Nil(t, result)
			assert.Error(t, err)
		})

		t.Run("Given template is successfully parsed with parameters", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingPromptFromPromptTemplateBeforeEach(t)

			agentId := uuid.New()
			parameters := map[string]any{"Name": "TestAgent"}
			template := "You are {{.Name}}, an expert system."

			suite.agentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertSystem).
				Return(&template, nil)

			result, err := suite.sut.ParseExpertSystem(agentId, parameters)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "You are TestAgent, an expert system.", *result)
		})
	})

	t.Run("ParseExpertContextFound", func(t *testing.T) {
		t.Parallel()

		t.Run("Given template requester returns error", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingPromptFromPromptTemplateBeforeEach(t)

			agentId := uuid.New()
			parameters := map[string]any{"Context": "some context"}

			suite.agentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertContextFound).
				Return(nil, assert.AnError)

			result, err := suite.sut.ParseExpertContextFound(agentId, parameters)

			assert.Nil(t, result)
			assert.Error(t, err)
		})

		t.Run("Given template is successfully parsed with parameters", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingPromptFromPromptTemplateBeforeEach(t)

			agentId := uuid.New()
			parameters := map[string]any{"Context": "relevant data"}
			template := "Context found: {{.Context}}"

			suite.agentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertContextFound).
				Return(&template, nil)

			result, err := suite.sut.ParseExpertContextFound(agentId, parameters)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "Context found: relevant data", *result)
		})
	})

	t.Run("ParseFacilitatorSystem", func(t *testing.T) {
		t.Parallel()

		t.Run("Given template with Go template syntax substitutes correctly", func(t *testing.T) {
			t.Parallel()

			suite := WhenParsingPromptFromPromptTemplateBeforeEach(t)

			agentId := uuid.New()
			parameters := map[string]any{
				"AgentName": "Facilitator",
				"Role":      "coordinator",
			}
			template := "{{.AgentName}} acts as a {{.Role}}."

			suite.agentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyFacilitatorSystem).
				Return(&template, nil)

			result, err := suite.sut.ParseFacilitatorSystem(agentId, parameters)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "Facilitator acts as a coordinator.", *result)
		})
	})
}
