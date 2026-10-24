package helpers_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenNormalizingCallerId(t *testing.T) {
	t.Parallel()

	t.Run("Given a padded mixed-case caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should trim and lowercase it", func(t *testing.T) {
			t.Parallel()

			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())

			normalized := helpers.NormalizeCallerId(" " + strings.ToUpper(callerId) + " ")

			assert.Equal(t, callerId, normalized)
		})
	})

	t.Run("Given an opaque caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "apikey:7f3c", helpers.NormalizeCallerId("ApiKey:7F3c"))
		})
	})

	t.Run("Given a blank caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty string", func(t *testing.T) {
			t.Parallel()

			assert.Empty(t, helpers.NormalizeCallerId("   "))
		})
	})
}
