package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenImportingAgentPreviousSignerTestingSuite struct {
	sut *services.AgentPreviousSignerImporter
}

func WhenImportingAgentPreviousSignerBeforeEach(t *testing.T) *WhenImportingAgentPreviousSignerTestingSuite {
	sut := services.NewAgentPreviousSignerImporter()

	return &WhenImportingAgentPreviousSignerTestingSuite{
		sut: sut,
	}
}

func TestWhenImportingAgentPreviousSigner(t *testing.T) {
	t.Parallel()

	t.Run("Given import is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return not implemented error", func(t *testing.T) {
			t.Parallel()
			basics := WhenImportingAgentPreviousSignerBeforeEach(t)

			agent := model_builders.NewAgentBuilder().Build()

			err := basics.sut.Import("gateway.example.com", agent)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "not implemented")
		})
	})
}
