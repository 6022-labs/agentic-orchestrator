package helpers_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_blockchain/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenEncodingEnsNames(t *testing.T) {
	t.Parallel()

	t.Run("Given a well-known ENS name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should compute the standard namehash", func(t *testing.T) {
			t.Parallel()

			node := helpers.Namehash("eth")
			assert.Equal(t, "93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae", hex.EncodeToString(node[:]))

			empty := helpers.Namehash("")
			assert.Equal(t, strings.Repeat("00", 32), hex.EncodeToString(empty[:]))
		})
	})

	t.Run("Given a multi-label name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should dns-encode it with length prefixes and a null terminator", func(t *testing.T) {
			t.Parallel()

			encoded, err := helpers.DnsEncode("agent.collection.eth")

			assert.NoError(t, err)
			assert.Equal(t, append(append(append([]byte{5}, []byte("agent")...), append([]byte{10}, []byte("collection")...)...), append([]byte{3}, []byte("eth")...)...), encoded[:len(encoded)-1])
			assert.Equal(t, byte(0), encoded[len(encoded)-1])
		})

		t.Run("Should reject labels longer than 63 bytes", func(t *testing.T) {
			t.Parallel()

			longLabel := make([]byte, 64)
			for i := range longLabel {
				longLabel[i] = 'a'
			}

			_, err := helpers.DnsEncode(string(longLabel) + ".eth")

			assert.Error(t, err)
		})
	})
}
