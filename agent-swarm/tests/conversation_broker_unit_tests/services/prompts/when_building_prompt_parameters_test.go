package prompts_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/prompts"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenBuildingPromptParametersTestingSuite struct {
	sut *prompts.PromptParametersBuilder

	mockRegistryAgentResolver *interfaces_mocks.MockRegistryAgentResolverInterface
}

func WhenBuildingPromptParametersBeforeEach(t *testing.T) *WhenBuildingPromptParametersTestingSuite {
	mockController := gomock.NewController(t)

	mockRegistryAgentResolver := interfaces_mocks.NewMockRegistryAgentResolverInterface(mockController)

	return &WhenBuildingPromptParametersTestingSuite{
		sut:                       prompts.NewPromptParametersBuilder(mockRegistryAgentResolver),
		mockRegistryAgentResolver: mockRegistryAgentResolver,
	}
}

// renderWith is how a stored template reaches these parameters, so it is how the test reads them.
func renderWith(t *testing.T, body string, parameters map[string]any) string {
	t.Helper()

	tmpl, err := template.New("test").Parse(body)
	assert.NoError(t, err)

	var out bytes.Buffer
	assert.NoError(t, tmpl.Execute(&out, parameters))
	return out.String()
}

func TestWhenBuildingPromptParameters(t *testing.T) {
	t.Parallel()

	t.Run("Given a member the registry knows", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain("hermes.eth").
				Return(model_builders.NewRegistryAgentBuilder().
					WithName("hermes").
					WithRole("clone").
					WithAddress(models.EvmAddressType, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").
					WithAddress("btc", "bc1qexample").
					Build(), nil).
				Times(1)
		}

		t.Run("Should expose every published address, keyed by chain", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)
			initSuite(suite)

			parameters, err := suite.sut.
				WithParam(
					[]string{"Members"},
					[]string{"hermes.eth"},
					lo.ToPtr(prompts.PromptParameterTypePotentialAgentArray),
				).
				Build()

			assert.NoError(t, err)
			rendered := renderWith(
				t,
				`{{range .Members}}{{index .Addresses "evm"}}|{{index .Addresses "btc"}}{{end}}`,
				parameters,
			)
			assert.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266|bc1qexample", rendered)
		})

		t.Run("Should expose the role through attributes rather than a field of its own", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)
			initSuite(suite)

			parameters, err := suite.sut.
				WithParam(
					[]string{"Members"},
					[]string{"hermes.eth"},
					lo.ToPtr(prompts.PromptParameterTypePotentialAgentArray),
				).
				Build()

			assert.NoError(t, err)
			rendered := renderWith(t, `{{range .Members}}{{.Name}}={{index .Attributes "role"}}{{end}}`, parameters)
			assert.Equal(t, "hermes=clone", rendered)
		})
	})

	t.Run("Given a member the registry does not know", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenBuildingPromptParametersTestingSuite) {
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain("pierre").
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should keep the raw name and mark it as not an agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)
			initSuite(suite)

			parameters, err := suite.sut.
				WithParam(
					[]string{"Members"},
					[]string{"pierre"},
					lo.ToPtr(prompts.PromptParameterTypePotentialAgentArray),
				).
				Build()

			assert.NoError(t, err)
			rendered := renderWith(t, `{{range .Members}}{{.Name}}:{{.IsAgent}}{{end}}`, parameters)
			assert.Equal(t, "pierre:false", rendered)
		})
	})

	t.Run("Given a plain parameter with no type directive", func(t *testing.T) {
		t.Parallel()

		t.Run("Should pass it through untouched", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingPromptParametersBeforeEach(t)
			suite.mockRegistryAgentResolver.EXPECT().ResolveByEnsDomain(gomock.Any()).Times(0)

			parameters, err := suite.sut.
				WithParam([]string{"ChannelDescription"}, "a test channel", nil).
				Build()

			assert.NoError(t, err)
			assert.Equal(t, "a test channel", renderWith(t, `{{.ChannelDescription}}`, parameters))
		})
	})
}
