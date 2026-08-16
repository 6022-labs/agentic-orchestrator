package helpers_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenComputingNamehash(t *testing.T) {
	t.Parallel()

	t.Run("Given a well-known ENS name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should compute the standard namehash", func(t *testing.T) {
			t.Parallel()

			node := helpers.Namehash("eth")
			assert.Equal(
				t,
				"93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae",
				hex.EncodeToString(node[:]),
			)
		})
	})

	t.Run("Given the 6022 base zone", func(t *testing.T) {
		t.Parallel()

		t.Run("Should match the registry's deployed parent node", func(t *testing.T) {
			t.Parallel()

			node := helpers.Namehash("6022.eth")
			assert.Equal(
				t,
				"9ec94752dcac587079728bde5f535457f7aa0bbe0b088555bca6d40cfcc78197",
				hex.EncodeToString(node[:]),
			)
		})
	})

	t.Run("Given an empty name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the zero node", func(t *testing.T) {
			t.Parallel()

			node := helpers.Namehash("")
			assert.Equal(t, strings.Repeat("00", 32), hex.EncodeToString(node[:]))
		})
	})
}
