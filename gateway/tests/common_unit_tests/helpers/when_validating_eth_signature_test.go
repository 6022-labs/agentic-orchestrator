package helpers_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingEthSignature(t *testing.T) {
	t.Parallel()

	t.Run("Given a 65-byte 0x signature", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			assert.True(t, helpers.IsEthSignature("0x"+strings.Repeat("ab", 65)))
		})
	})

	t.Run("Given a signature with the wrong length or no prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsEthSignature("0x"+strings.Repeat("ab", 64)))
			assert.False(t, helpers.IsEthSignature(strings.Repeat("ab", 65)))
			assert.False(t, helpers.IsEthSignature("0xnot-hex"))
		})
	})
}
