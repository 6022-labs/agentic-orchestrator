package prompt_parameters_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/compute/prompt_parameters"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenBuildingPromptParametersTestingSuite struct {
	sut                        interfaces.PromptParametersBuilderInterface
	mockAgentRegistryRequester *interfaces_mocks.MockAgentRegistryRequesterInterface
}

func WhenBuildingPromptParametersBeforeEach(t *testing.T) *WhenBuildingPromptParametersTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)

	sut := prompt_parameters.NewPromptParametersBuilder(mockAgentRegistryRequester)

	return &WhenBuildingPromptParametersTestingSuite{
		sut:                        sut,
		mockAgentRegistryRequester: mockAgentRegistryRequester,
	}
}

func TestWhenBuildingPromptParameters(t *testing.T) {
	t.Parallel()

	t.Run("Given no parameters are set", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty map", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			result, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Empty(t, result)
		})
	})

	t.Run("Given base parameters are set", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return base parameters", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			base := map[string]any{
				"key1": "value1",
				"key2": "value2",
			}

			result, err := suite.sut.WithBase(base).Build()

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "value1", result["key1"])
			assert.Equal(t, "value2", result["key2"])
		})
	})

	t.Run("Given simple parameters are added", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return parameters with simple values", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			result, err := suite.sut.
				WithParam([]string{"SimpleKey"}, "simpleValue", nil).
				Build()

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "simpleValue", result["SimpleKey"])
		})
	})

	t.Run("Given nested parameters are added", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return parameters with nested values", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			result, err := suite.sut.
				WithParam([]string{"Parent", "Child"}, "nestedValue", nil).
				Build()

			assert.NoError(t, err)
			assert.NotNil(t, result)

			parent, ok := result["Parent"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, "nestedValue", parent["Child"])
		})
	})

	t.Run("Given a parameter with agent type is added", func(t *testing.T) {
		t.Parallel()

		t.Run("And agent is found", func(t *testing.T) {
			t.Parallel()

			agentResponse := model_builders.NewAgentResponseBuilder().
				WithEnsDomain("test-agent.protocol6022.eth").
				Build()

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetMintedAgentByEnsDomain("test-agent.protocol6022.eth").
					Return(agentResponse, nil)
			}

			t.Run("Should resolve agent parameter", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				agentType := interfaces.PromptParameterTypeAgent
				result, err := suite.sut.
					WithParam([]string{"Agent"}, "test-agent.protocol6022.eth", &agentType).
					Build()

				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, "test-agent.protocol6022.eth", result["Agent"])
			})
		})

		t.Run("And agent is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetMintedAgentByEnsDomain("unknown-agent.protocol6022.eth").
					Return(nil, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				agentType := interfaces.PromptParameterTypeAgent
				result, err := suite.sut.
					WithParam([]string{"Agent"}, "unknown-agent.protocol6022.eth", &agentType).
					Build()

				assert.Error(t, err)
				assert.Nil(t, result)
			})
		})

		t.Run("And agent registry returns error", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetMintedAgentByEnsDomain("error-agent.protocol6022.eth").
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				agentType := interfaces.PromptParameterTypeAgent
				result, err := suite.sut.
					WithParam([]string{"Agent"}, "error-agent.protocol6022.eth", &agentType).
					Build()

				assert.Error(t, err)
				assert.Nil(t, result)
			})
		})
	})

	t.Run("Given a parameter with agent array type is added", func(t *testing.T) {
		t.Parallel()

		t.Run("And all agents are found", func(t *testing.T) {
			t.Parallel()

			agent1 := model_builders.NewAgentResponseBuilder().
				WithEnsDomain("agent1.protocol6022.eth").
				Build()
			agent2 := model_builders.NewAgentResponseBuilder().
				WithEnsDomain("agent2.protocol6022.eth").
				Build()

			domains := []string{"agent1.protocol6022.eth", "agent2.protocol6022.eth"}

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetAgentsByEnsDomains(domains).
					Return([]*responses.AgentResponse{agent1, agent2}, nil)
			}

			t.Run("Should resolve agent array parameter", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				agentArrayType := interfaces.PromptParameterTypeAgentArray
				result, err := suite.sut.
					WithParam([]string{"Agents"}, domains, &agentArrayType).
					Build()

				assert.NoError(t, err)
				assert.NotNil(t, result)
			})
		})

		t.Run("And not all agents are found", func(t *testing.T) {
			t.Parallel()

			agent1 := model_builders.NewAgentResponseBuilder().
				WithEnsDomain("agent1.protocol6022.eth").
				Build()

			domains := []string{"agent1.protocol6022.eth", "agent2.protocol6022.eth"}

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetAgentsByEnsDomains(domains).
					Return([]*responses.AgentResponse{agent1}, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				agentArrayType := interfaces.PromptParameterTypeAgentArray
				result, err := suite.sut.
					WithParam([]string{"Agents"}, domains, &agentArrayType).
					Build()

				assert.Error(t, err)
				assert.Nil(t, result)
			})
		})
	})

	t.Run("Given a parameter with potential agent type is added", func(t *testing.T) {
		t.Parallel()

		t.Run("And agent is found", func(t *testing.T) {
			t.Parallel()

			agentResponse := model_builders.NewAgentResponseBuilder().
				WithEnsDomain("potential.protocol6022.eth").
				Build()

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetMintedAgentByEnsDomain("potential.protocol6022.eth").
					Return(agentResponse, nil)
			}

			t.Run("Should resolve potential agent parameter", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				potentialAgentType := interfaces.PromptParameterTypePotentialAgent
				result, err := suite.sut.
					WithParam([]string{"Author"}, "potential.protocol6022.eth", &potentialAgentType).
					Build()

				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, "potential.protocol6022.eth", result["Author"])
			})
		})

		t.Run("And agent is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetMintedAgentByEnsDomain("not-agent.protocol6022.eth").
					Return(nil, nil)
			}

			t.Run("Should still succeed with potential agent parameter", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				potentialAgentType := interfaces.PromptParameterTypePotentialAgent
				result, err := suite.sut.
					WithParam([]string{"Author"}, "not-agent.protocol6022.eth", &potentialAgentType).
					Build()

				assert.NoError(t, err)
				assert.NotNil(t, result)
			})
		})
	})

	t.Run("Given base parameters contain type directives", func(t *testing.T) {
		t.Parallel()

		t.Run("And type is agent", func(t *testing.T) {
			t.Parallel()

			agentResponse := model_builders.NewAgentResponseBuilder().
				WithEnsDomain("base-agent.protocol6022.eth").
				Build()

			initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
				suite.mockAgentRegistryRequester.EXPECT().
					GetMintedAgentByEnsDomain("base-agent.protocol6022.eth").
					Return(agentResponse, nil)
			}

			t.Run("Should resolve type directives from base", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingPromptParametersBeforeEach(t)
				initSuite(suite)

				base := map[string]any{
					"BaseAgent":      "base-agent.protocol6022.eth",
					"BaseAgent.type": "agent",
				}

				result, err := suite.sut.WithBase(base).Build()

				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, "base-agent.protocol6022.eth", result["BaseAgent"])
				assert.Nil(t, result["BaseAgent.type"])
			})
		})
	})

	t.Run("Given empty path is provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should ignore the parameter", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			result, err := suite.sut.
				WithParam([]string{}, "ignoredValue", nil).
				Build()

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Empty(t, result)
		})
	})
}
