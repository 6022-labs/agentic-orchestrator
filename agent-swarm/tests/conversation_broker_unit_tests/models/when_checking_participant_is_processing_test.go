package models_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWhenCheckingParticipantIsProcessing(t *testing.T) {
	t.Parallel()

	t.Run("Given no lease", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be processing", func(t *testing.T) {
			t.Parallel()

			assert.False(t, model_builders.NewAgentParticipantBuilder().Build().IsProcessing())
		})
	})

	t.Run("Given a running lease", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be processing", func(t *testing.T) {
			t.Parallel()

			p := model_builders.NewAgentParticipantBuilder().WithLease(uuid.New(), time.Now().Add(time.Minute)).Build()

			assert.True(t, p.IsProcessing())
		})
	})

	t.Run("Given an expired lease", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be processing", func(t *testing.T) {
			t.Parallel()

			p := model_builders.NewAgentParticipantBuilder().WithLease(uuid.New(), time.Now().Add(-time.Minute)).Build()

			assert.False(t, p.IsProcessing())
		})
	})
}
