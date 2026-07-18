package helpers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingEthAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given a checksummed 0x address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			assert.True(t, helpers.IsEthAddress("0xd389d997Df5654b87914492E1e0C6861bfe835DE"))
		})
	})

	t.Run("Given an address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsEthAddress("d389d997Df5654b87914492E1e0C6861bfe835DE"))
		})
	})

	t.Run("Given a malformed address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			assert.False(t, helpers.IsEthAddress("0x1234"))
			assert.False(t, helpers.IsEthAddress(""))
			assert.False(t, helpers.IsEthAddress("0xZZ89d997Df5654b87914492E1e0C6861bfe835DE"))
		})
	})
}
