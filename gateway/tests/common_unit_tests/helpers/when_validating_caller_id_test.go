package helpers_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingCallerId(t *testing.T) {
	t.Parallel()

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsCallerId(""))
		})
	})

	t.Run("Given a blank caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsCallerId("   "))
		})
	})

	t.Run("Given a caller id longer than the maximum length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsCallerId(strings.Repeat("k", helpers.CallerIdMaxLength+1)))
		})
	})

	t.Run("Given a multibyte caller id longer than the maximum length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsCallerId(strings.Repeat("é", helpers.CallerIdMaxLength+1)))
		})
	})

	t.Run("Given a caller id of exactly the maximum length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			assert.True(t, helpers.IsCallerId(strings.Repeat("é", helpers.CallerIdMaxLength)))
		})
	})

	t.Run("Given an opaque caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			assert.True(t, helpers.IsCallerId("ApiKey:7F3c"))
		})
	})
}
