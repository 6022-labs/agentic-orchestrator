package helpers_test

import (
	"context"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenCarryingCallerId(t *testing.T) {
	t.Parallel()

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report no caller", func(t *testing.T) {
			t.Parallel()

			_, ok := helpers.CallerIdFromContext(context.Background())

			assert.False(t, ok)
		})
	})

	t.Run("Given a caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should read it back", func(t *testing.T) {
			t.Parallel()

			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)

			carried, ok := helpers.CallerIdFromContext(ctx)

			assert.True(t, ok)
			assert.Equal(t, callerId, carried)
		})
	})
}
