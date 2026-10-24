package helpers_test

import (
	"context"
	"strings"
	"testing"

	apperrors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenRequiringCallerId(t *testing.T) {
	t.Parallel()

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			_, err := helpers.RequireCallerId(context.Background())

			var unauthorized *apperrors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorized)
			assert.Equal(t, "missing_caller_id", unauthorized.Code)
		})
	})

	t.Run("Given a caller id longer than the maximum length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			ctx := helpers.ContextWithCallerId(
				context.Background(),
				strings.Repeat("k", helpers.CallerIdMaxLength+1),
			)

			_, err := helpers.RequireCallerId(ctx)

			var ve *apperrors.ValidationError
			assert.ErrorAs(t, err, &ve)
			assert.Equal(t, "callerId", ve.Field)
		})
	})

	t.Run("Given a caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it", func(t *testing.T) {
			t.Parallel()

			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)

			required, err := helpers.RequireCallerId(ctx)

			assert.NoError(t, err)
			assert.Equal(t, callerId, required)
		})
	})
}
