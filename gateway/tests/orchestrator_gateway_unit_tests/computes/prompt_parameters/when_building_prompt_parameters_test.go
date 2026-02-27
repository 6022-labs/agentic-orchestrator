package prompt_parameters_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/compute/prompt_parameters"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenBuildingPromptParametersTestingSuite struct {
	sut *prompt_parameters.PromptParametersBuilder

	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface
}

func WhenBuildingPromptParametersBeforeEach(t *testing.T) *WhenBuildingPromptParametersTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := prompt_parameters.NewPromptParametersBuilder(agentRepositoryMock)

	return &WhenBuildingPromptParametersTestingSuite{
		sut: sut,

		agentRepositoryMock: agentRepositoryMock,
	}
}

func TestWhenBuildingPromptParameters(t *testing.T) {
	t.Parallel()

	toJsonMap := func(val any) map[string]any {
		data, _ := json.Marshal(val)
		var result map[string]any
		json.Unmarshal(data, &result)
		return result
	}

	toJsonArray := func(val any) []map[string]any {
		data, _ := json.Marshal(val)
		var result []map[string]any
		json.Unmarshal(data, &result)
		return result
	}

	t.Run("Given no hints are present", func(t *testing.T) {
		t.Parallel()

		t.Run("Should keep base map", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			base := map[string]any{
				"foo": "bar",
				"num": 42,
			}

			suite.sut.WithBase(base)
			out, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.Equal(t, "bar", out["foo"])
			assert.Equal(t, 42, out["num"])
		})
	})

	t.Run("Given type hint is agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should resolve agent by ENS", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			agent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.EXPECT().GetByEnsDomain(agent.EnsDomain).Return(agent, nil)

			base := map[string]any{
				"target":      agent.EnsDomain,
				"target.type": "agent",
				"keep":        "me",
			}

			suite.sut.WithBase(base)
			out, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.Equal(t, "me", out["keep"])
			_, hasType := out["target.type"]
			assert.False(t, hasType, "hint key should be removed")

			asJSON := toJsonMap(out["target"])
			assert.Equal(t, agent.Name, asJSON["name"])
			assert.Equal(t, *agent.Owner, asJSON["owner"])
			assert.Equal(t, agent.Role.String(), asJSON["role"])
			assert.Equal(t, agent.EnsDomain, asJSON["ensDomain"])
			attrs, ok := asJSON["attributes"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, toJsonMap(agent.Attributes), attrs)
			assert.Equal(t, agent.WalletAddress, asJSON["walletAddress"])
		})
	})

	t.Run("Given setting additional WithParam", func(t *testing.T) {
		t.Parallel()

		t.Run("Should set value at root for single-segment path", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			suite.sut.WithParam([]string{"alpha"}, 123, nil)
			out, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.Equal(t, 123, out["alpha"])
		})

		t.Run("Should create nested maps for multi-segment path", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			suite.sut.WithParam([]string{"object", "inner", "leaf"}, "value", nil)
			out, err := suite.sut.Build()

			assert.NoError(t, err)

			lvl1, ok := out["object"].(map[string]any)
			assert.True(t, ok)
			lvl2, ok := lvl1["inner"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, "value", lvl2["leaf"])
		})

		t.Run("Should overwrite non-map intermediate value when creating nested path", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			// Start with a non-map at key "obj"
			suite.sut.WithBase(map[string]any{
				"obj": "was-string",
			})

			// This should replace "obj" with a map containing {"inner": 7}
			suite.sut.WithParam([]string{"obj", "inner"}, 7, nil)
			out, err := suite.sut.Build()

			assert.NoError(t, err)

			obj, ok := out["obj"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, 7, obj["inner"])
		})

		t.Run("Should ignore call when path is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			suite.sut.WithBase(map[string]any{"keep": "me"})
			suite.sut.WithParam([]string{}, "ignored", nil)
			out, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.Equal(t, "me", out["keep"])
			// Ensure nothing unexpected was added at root
			_, hasEmpty := out[""]
			assert.False(t, hasEmpty)
		})

		t.Run("Should add type hint when provided and resolve value, then remove hint on build", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			pt := interfaces.PromptParameterTypeAgent
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.EXPECT().GetByEnsDomain(agent.EnsDomain).Return(agent, nil)

			suite.sut.WithParam([]string{"target"}, agent.EnsDomain, &pt)
			out, err := suite.sut.Build()

			assert.NoError(t, err)
			asJSON := toJsonMap(out["target"])
			assert.Equal(t, agent.Name, asJSON["name"])
			assert.Equal(t, *agent.Owner, asJSON["owner"])
			assert.Equal(t, agent.Role.String(), asJSON["role"])
			assert.Equal(t, agent.EnsDomain, asJSON["ensDomain"])
			attrs, ok := asJSON["attributes"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, toJsonMap(agent.Attributes), attrs)
			assert.Equal(t, agent.WalletAddress, asJSON["walletAddress"])
			// The .type hint should have been removed during build processing
			_, hasType := out["target.type"]
			assert.False(t, hasType)
		})

		t.Run("Should not add type hint when nil or empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			// nil type
			suite.sut.WithParam([]string{"nilType"}, "x", nil)

			// empty type
			empty := interfaces.PromptParameterType("")
			suite.sut.WithParam([]string{"emptyType"}, "y", &empty)

			out, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.Equal(t, "x", out["nilType"])
			assert.Equal(t, "y", out["emptyType"])
			_, hasNilType := out["nilType.type"]
			_, hasEmptyType := out["emptyType.type"]
			assert.False(t, hasNilType)
			assert.False(t, hasEmptyType)
		})

		t.Run("Should support chaining multiple WithParam calls", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			suite.sut.
				WithParam([]string{"a"}, 1, nil).
				WithParam([]string{"b", "c"}, "z", nil)

			out, err := suite.sut.Build()

			assert.NoError(t, err)
			assert.Equal(t, 1, out["a"])
			b, ok := out["b"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, "z", b["c"])
		})
	})

	t.Run("Given agent not found", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			ens := "missing.eth"
			suite.agentRepositoryMock.EXPECT().GetByEnsDomain(ens).Return(nil, gorm.ErrRecordNotFound)

			suite.sut.WithBase(map[string]any{
				"target":      ens,
				"target.type": "agent",
			})

			_, err := suite.sut.Build()
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while requesting repository", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {

			suite := WhenBuildingPromptParametersBeforeEach(t)
			ens := "err.eth"
			suite.agentRepositoryMock.EXPECT().GetByEnsDomain(ens).Return(nil, errors.New("db down"))

			suite.sut.WithBase(map[string]any{
				"target":      ens,
				"target.type": "agent",
			})

			_, err := suite.sut.Build()
			assert.Error(t, err)
		})
	})

	t.Run("Given wrong value type for agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {

			suite := WhenBuildingPromptParametersBeforeEach(t)
			suite.sut.WithBase(map[string]any{
				"target":      123,
				"target.type": "agent",
			})

			_, err := suite.sut.Build()
			assert.Error(t, err)
		})
	})

	t.Run("Given type hint is agent[] and de-duplicate", func(t *testing.T) {
		t.Parallel()

		t.Run("Should resolve agents by ENS", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			a1 := model_builders.NewAgentBuilder().Build()
			a2 := model_builders.NewAgentBuilder().Build()

			// Expect deduped call order-preserving
			suite.agentRepositoryMock.EXPECT().GetByEnsDomains([]string{a1.EnsDomain, a2.EnsDomain}).Return([]*models.Agent{a1, a2}, nil)

			base := map[string]any{
				"targets":      []string{a1.EnsDomain, a1.EnsDomain, a2.EnsDomain},
				"targets.type": "agents", // alias of agent[]
			}

			suite.sut.WithBase(base)
			out, err := suite.sut.Build()
			assert.NoError(t, err)

			_, hasType := out["targets.type"]
			assert.False(t, hasType)

			arr := toJsonArray(out["targets"])

			assert.Equal(t, a1.Name, arr[0]["name"])
			assert.Equal(t, a1.Role.String(), arr[0]["role"])
			assert.Equal(t, a2.Name, arr[1]["name"])
			assert.Equal(t, a2.Role.String(), arr[1]["role"])
		})
	})

	t.Run("Given not all agents found", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error for agent[]", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			firstAgent := model_builders.NewAgentBuilder().Build()

			suite.agentRepositoryMock.EXPECT().GetByEnsDomains([]string{"alice.eth", "bob.eth"}).Return([]*models.Agent{firstAgent}, nil)

			suite.sut.WithBase(map[string]any{
				"targets":      []string{"alice.eth", "bob.eth"},
				"targets.type": "agent[]",
			})

			_, err := suite.sut.Build()
			assert.Error(t, err)
		})
	})

	t.Run("Given wrong value type for agent[]", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)
			suite.sut.WithBase(map[string]any{
				"targets":      "not-a-slice",
				"targets.type": "agent[]",
			})

			_, err := suite.sut.Build()
			assert.Error(t, err)
		})
	})

	t.Run("Given type hint is potential_agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should resolve agent by ENS and keep provided name", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.EXPECT().GetByEnsDomain(agent.EnsDomain).Return(agent, nil)

			suite.sut.WithBase(map[string]any{
				"pot":      agent.EnsDomain,
				"pot.type": "potential_agent",
			})

			out, err := suite.sut.Build()
			assert.NoError(t, err)

			pot := toJsonMap(out["pot"])
			assert.Equal(t, true, pot["isAgent"])
			assert.Equal(t, agent.Name, pot["name"])
			assert.Equal(t, agent.Role.String(), pot["role"])
			assert.Equal(t, agent.WalletAddress, pot["walletAddress"])
		})
	})

	t.Run("Given potential_agent is not found", func(t *testing.T) {
		t.Parallel()

		t.Run("Should set IsAgent to false", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			suite.agentRepositoryMock.EXPECT().GetByEnsDomain("ghost.eth").Return(nil, gorm.ErrRecordNotFound)

			suite.sut.WithBase(map[string]any{
				"pot":      "ghost.eth",
				"pot.type": "potential_agent",
			})

			out, err := suite.sut.Build()
			assert.NoError(t, err)

			pot := toJsonMap(out["pot"])
			assert.Equal(t, "ghost.eth", pot["name"])
			assert.Equal(t, false, pot["isAgent"])
			// Agent fields should be absent
			_, hasRole := pot["role"]
			assert.False(t, hasRole)
		})
	})

	t.Run("Given type hint is potential_agent[]", func(t *testing.T) {
		t.Parallel()

		t.Run("Should resolve mix of existing and missing", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)

			firstAgent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.EXPECT().GetByEnsDomains([]string{firstAgent.EnsDomain, "bob.eth"}).Return([]*models.Agent{firstAgent}, nil)

			suite.sut.WithBase(map[string]any{
				"pots":      []string{firstAgent.EnsDomain, "bob.eth"},
				"pots.type": "potential_agents", // alias
			})

			out, err := suite.sut.Build()
			assert.NoError(t, err)

			arr := toJsonArray(out["pots"])
			assert.Len(t, arr, 2)

			first := arr[0]
			second := arr[1]

			assert.Equal(t, true, first["isAgent"])
			assert.Equal(t, firstAgent.Name, first["name"])
			assert.Equal(t, firstAgent.Role.String(), first["role"])

			assert.Equal(t, nil, second["role"])
			assert.Equal(t, false, second["isAgent"])
			assert.Equal(t, "bob.eth", second["name"])
		})
	})
}
